package usable_accept_rate

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"emcontroller/auto-schedule/algorithms"
	applicationsgenerator "emcontroller/auto-schedule/experiments/applications-generator"
	asmodel "emcontroller/auto-schedule/model"
	"emcontroller/models"
)

const dataFileName string = "usable_acceptance_rate.csv"

// the data structure that will be collected in this experiment
type exptData struct {
	algorithmName string

	schedulingRequestCount int
	usableSolutionCount    int

	totalAppCount         int
	totalAcceptedAppCount int

	appCountPerPri         map[int]int
	acceptedAppCountPerPri map[int]int

	totalAppPriority         int
	totalAcceptedAppPriority int

	solutionUsableRate                float64
	appAcceptanceRate                 float64
	appPriorityWeightedAcceptanceRate float64
	appPerPriAcceptanceRate           map[int]float64
}

func TestExecute(t *testing.T) {
	var appNamePrefix string = "expt-app"
	var appCount int = 100
	var repeatCount int = 2 // We repeat this experiment for 10 times to reduce the impact from random factors.

	// all algorithms to be evaluated in experiment
	var algoNames []string = []string{algorithms.CompRandName, algorithms.BERandName, algorithms.AmagaName, algorithms.AmpgaName, algorithms.McssgaName}

	var results []exptData // used to save and output results
	for _, algoName := range algoNames {
		results = append(results, exptData{
			algorithmName: algoName, appCountPerPri: make(map[int]int), acceptedAppCountPerPri: make(map[int]int), appPerPriAcceptanceRate: make(map[int]float64),
		})
	}

	// We repeat experiment to reduce the impact from random factors. In every repeat, we generate different applications.
	for i := 0; i < repeatCount; i++ {
		apps, err := applicationsgenerator.MakeExperimentApps(appNamePrefix, appCount, false)
		if err != nil {
			t.Errorf("MakeExperimentApps error: %s", err.Error())
		}
		for j, algoName := range algoNames { // in one repeat, we use the same apps for all algorithm for comparison.
			t.Logf("Repeat %d, algorithm No. %d [%s]", i, j, algoName)

			acceptedApps, usable, err := schedulingRequest(algoName, apps)
			if err != nil {
				t.Errorf("schedulingRequest error: %s", err.Error())
			}

			// record results
			results[j].schedulingRequestCount++
			results[j].totalAppCount += len(apps)
			for _, app := range apps {
				results[j].totalAppPriority += app.Priority
			}
			appCountPerPri := getPerPriAppCount(apps)
			for pri := asmodel.MinPriority; pri <= asmodel.MaxPriority; pri++ {
				results[j].appCountPerPri[pri] += appCountPerPri[pri]
			}

			if usable {
				results[j].usableSolutionCount++
				results[j].totalAcceptedAppCount += len(acceptedApps)
				for _, acceptedApp := range acceptedApps {
					results[j].totalAcceptedAppPriority += acceptedApp.Priority
				}
				acceptedAppCountPerPri := getPerPriAcceptedAppCount(acceptedApps)
				for pri := asmodel.MinPriority; pri <= asmodel.MaxPriority; pri++ {
					results[j].acceptedAppCountPerPri[pri] += acceptedAppCountPerPri[pri]
				}
			}
		}
	}

	// calculate the rates in the results
	for i := 0; i < len(results); i++ {
		results[i].solutionUsableRate = float64(results[i].usableSolutionCount) / float64(results[i].schedulingRequestCount)
		results[i].appAcceptanceRate = float64(results[i].totalAcceptedAppCount) / float64(results[i].totalAppCount)
		results[i].appPriorityWeightedAcceptanceRate = float64(results[i].totalAcceptedAppPriority) / float64(results[i].totalAppPriority)
		for pri := asmodel.MinPriority; pri <= asmodel.MaxPriority; pri++ {
			results[i].appPerPriAcceptanceRate[pri] = float64(results[i].acceptedAppCountPerPri[pri]) / float64(results[i].appCountPerPri[pri])
		}
	}

	if err := writeCsvResults(results); err != nil {
		t.Errorf("writeCsvResults error: %s", err.Error())
	}

}

func schedulingRequest(algoName string, apps []models.K8sApp) ([]models.AppInfo, bool, error) {
	url := "http://localhost:20000/doNewAppGroup"

	reqBodyJson, err := json.Marshal(apps)
	if err != nil {
		outErr := fmt.Errorf("json.Marshal: %+v, error: %w", apps, err)
		return []models.AppInfo{}, false, outErr
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBodyJson))
	if err != nil {
		return []models.AppInfo{}, false, fmt.Errorf("url: %s, make request error: %w", url, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Mcm-Scheduling-Algorithm", algoName)
	req.Header.Set("Expected-Time-One-Cpu", "35")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []models.AppInfo{}, false, fmt.Errorf("url: %s, do request error: %w", url, err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []models.AppInfo{}, false, fmt.Errorf("url: %s, res.statusCode is %d, read res.Body error: %w", url, res.StatusCode, err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		if strings.Contains(string(body), "unusable solution") { // the scheduling scheme is unusable
			return []models.AppInfo{}, false, nil // return of unusable solution
		}
		return []models.AppInfo{}, false, fmt.Errorf("url: %s, res.statusCode is %d, res.Body is %s", url, res.StatusCode, string(body))
	}

	var acceptedApps []models.AppInfo
	if err := json.Unmarshal(body, &acceptedApps); err != nil {
		return []models.AppInfo{}, true, fmt.Errorf("url: %s, res.statusCode is %d, res.Body is %s, Unmarshal body error: %s", url, res.StatusCode, string(body), err.Error())
	}

	return acceptedApps, true, nil // return of usable solution
}

// get the number of applications with each priority
func getPerPriAppCount(apps []models.K8sApp) map[int]int {
	var perPriAppCount map[int]int = make(map[int]int)

	for _, app := range apps {
		perPriAppCount[app.Priority]++
	}

	return perPriAppCount
}

// get the number of accepted applications with each priority
func getPerPriAcceptedAppCount(acceptedApps []models.AppInfo) map[int]int {
	var perPriAcceptedAppCount map[int]int = make(map[int]int)

	for _, acceptedApp := range acceptedApps {
		perPriAcceptedAppCount[acceptedApp.Priority]++
	}

	return perPriAcceptedAppCount
}

// function to write data into a csv file.
func writeCsvResults(results []exptData) error {

	var csvContent [][]string

	var header []string = []string{
		"Algorithm Name",
		"Scheduling Request Count",
		"Usable Solution Count",
		"Total App Count",
		"Total Accepted App Count",
		"Total App Priority",
		"Total Accepted App Priority",
		"Solution Usable Rate",
		"App Acceptance Rate",
		"App Priority Weighted Acceptance Rate",
	}
	for pri := asmodel.MinPriority; pri <= asmodel.MaxPriority; pri++ {
		header = append(header, fmt.Sprintf("Priority-%d App Acceptance Rate", pri))
	}
	csvContent = append(csvContent, header)

	for _, result := range results {
		var line []string = []string{
			result.algorithmName,
			fmt.Sprintf("%d", result.schedulingRequestCount),
			fmt.Sprintf("%d", result.usableSolutionCount),
			fmt.Sprintf("%d", result.totalAppCount),
			fmt.Sprintf("%d", result.totalAcceptedAppCount),
			fmt.Sprintf("%d", result.totalAppPriority),
			fmt.Sprintf("%d", result.totalAcceptedAppPriority),
			fmt.Sprintf("%g", result.solutionUsableRate),
			fmt.Sprintf("%g", result.appAcceptanceRate),
			fmt.Sprintf("%g", result.appPriorityWeightedAcceptanceRate),
		}
		for pri := asmodel.MinPriority; pri <= asmodel.MaxPriority; pri++ {
			line = append(line, fmt.Sprintf("%g", result.appPerPriAcceptanceRate[pri]))
		}
		csvContent = append(csvContent, line)
	}

	return writeCsvFile(dataFileName, csvContent)
}

func writeCsvFile(fileName string, csvContent [][]string) error {
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("create file %s, error: %w", fileName, err)
	}
	w := csv.NewWriter(f)
	defer w.Flush()

	for _, record := range csvContent {
		if err := w.Write(record); err != nil {
			return fmt.Errorf("write record %v, error: %s", record, err.Error())
		}
	}

	return nil
}
