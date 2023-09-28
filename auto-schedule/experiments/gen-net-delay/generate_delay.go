package gen_net_delay

import (
	"fmt"
	"sync"

	"github.com/astaxie/beego"

	"emcontroller/models"
)

// on each Proxmox server, iface "vmbr0" is a virtual port bridged on 2 physical ports "ens255f0" and "ens4f1". If we only add delay on the virtual port "vmbr0", the VMs inside the servers will not be affected, so we should add delay to the physical port. Now, for every server we are using "ens4f1", so we add delay to this port.
const netInterfaceName string = "ens4f1"

// generate network delay for clouds. The input map's keys are the name of clouds, and the input map's values are the delay to add to clouds.
func GenCloudsDelay(delays map[string]string) error {
	var errs []error

	// the slice in golang is not safe for concurrent read/write
	var errsMu sync.Mutex

	// delay every cloud in parallel
	var wg sync.WaitGroup

	for cloudName, delay := range delays {
		wg.Add(1)
		go func(cn, d string) {
			defer wg.Done()
			err := delayOneCloud(cn, d)
			if err != nil {
				outErr := fmt.Errorf("generate delay [%s] for cloud [%s], error %w.", d, cn, err)
				beego.Error(outErr)
				errsMu.Lock()
				errs = append(errs, outErr)
				errsMu.Unlock()
			} else {
				beego.Info(fmt.Sprintf("Successfully add delay [%s] for cloud [%s].", d, cn))
			}
		}(cloudName, delay)
	}
	wg.Wait()

	if len(errs) != 0 {
		sumErr := models.HandleErrSlice(errs)
		outErr := fmt.Errorf("Generate clouds network delay, Error: %w", sumErr)
		beego.Error(outErr)
		return outErr
	}

	return nil
}

// clear delay on all clouds
func ClearAllDelay() error {
	var errs []error

	// the slice in golang is not safe for concurrent read/write
	var errsMu sync.Mutex

	// delay every cloud in parallel
	var wg sync.WaitGroup

	for cloudName, _ := range models.Clouds {
		if _, ok := models.Clouds[cloudName].(*models.Proxmox); !ok {
			beego.Info(fmt.Sprintf("ClearAllDelay Skip cloud %s, because its type is not %s.", cloudName, models.ProxmoxIaas))
			continue
		}
		wg.Add(1)
		go func(cn string) {
			defer wg.Done()
			err := clearDelayOneCloud(cn)
			if err != nil {
				outErr := fmt.Errorf("clear for cloud [%s], error %w.", cn, err)
				beego.Error(outErr)
				errsMu.Lock()
				errs = append(errs, outErr)
				errsMu.Unlock()
			} else {
				beego.Info(fmt.Sprintf("Successfully clear delay for cloud [%s].", cn))
			}
		}(cloudName)
	}
	wg.Wait()

	if len(errs) != 0 {
		sumErr := models.HandleErrSlice(errs)
		outErr := fmt.Errorf("Clear all clouds network delay, Error: %w", sumErr)
		beego.Error(outErr)
		return outErr
	}

	return nil
}

func delayOneCloud(cloudName, delay string) error {
	cloud, exist := models.Clouds[cloudName]
	if !exist {
		return fmt.Errorf("delayOneCloud, cloud name [%s] not found", cloudName)
	}

	var sshUser, sshPwd, sshIp string
	switch realTypeCloud := cloud.(type) {
	case *models.Proxmox:
		sshUser = realTypeCloud.ProxmoxUser
		sshPwd = realTypeCloud.ProxmoxPassword
		sshIp = realTypeCloud.IP
	default:
		return fmt.Errorf("delayOneCloud, the type of cloud [%s] is [%s], which does not support to add network delay", cloudName, cloud.ShowType())
	}

	// clear old delay before adding new delay
	if err := clearDelay(sshUser, sshPwd, sshIp); err != nil {
		return fmt.Errorf("delayOneCloud, add delay to cloud name [%s], clear old delay error: [%s]", cloudName, err.Error())
	}

	return setDelay(sshUser, sshPwd, sshIp, delay)
}

func clearDelayOneCloud(cloudName string) error {
	cloud, exist := models.Clouds[cloudName]
	if !exist {
		return fmt.Errorf("delayOneCloud, cloud name [%s] not found", cloudName)
	}

	var sshUser, sshPwd, sshIp string
	switch realTypeCloud := cloud.(type) {
	case *models.Proxmox:
		sshUser = realTypeCloud.ProxmoxUser
		sshPwd = realTypeCloud.ProxmoxPassword
		sshIp = realTypeCloud.IP
	default:
		return fmt.Errorf("delayOneCloud, the type of cloud [%s] is [%s], which does not support to add network delay", cloudName, cloud.ShowType())
	}

	return clearDelay(sshUser, sshPwd, sshIp)
}

// ssh to a server/vm to add network delay using TC (traffic control)
func setDelay(sshUser, sshPwd, sshIp, delay string) error {
	clearDelayCmd := fmt.Sprintf("tc qdisc add dev %s root netem delay %s", netInterfaceName, delay)

	beego.Info(fmt.Sprintf("SSH to the IP [%s] to run command [%s] to set network delay.", sshIp, clearDelayCmd))
	sshClient, err := models.SshClientWithPasswd(sshUser, sshPwd, sshIp, models.SshPort)
	if err != nil {
		return fmt.Errorf("setDelay, create SshClientWithPasswd for ip [%s], error: %w", sshIp, err)
	}
	defer sshClient.Close()

	output, err := models.SshOneCommand(sshClient, clearDelayCmd)
	if err != nil {
		fmt.Errorf("setDelay, Execute command [%s] on IP [%s] error: [%s].", clearDelayCmd, sshIp, err.Error())
	}
	beego.Info(fmt.Sprintf("Set network delay [%s] on IP [%s], output: [%s].", delay, sshIp, string(output)))
	return nil
}

// ssh to a server/vm to clear network delay added by TC (traffic control)
func clearDelay(sshUser, sshPwd, sshIp string) error {
	clearDelayCmd := fmt.Sprintf("tc qdisc del dev %s root", netInterfaceName)

	beego.Info(fmt.Sprintf("SSH to the IP [%s] to run command [%s] to clear all added network delay.", sshIp, clearDelayCmd))
	sshClient, err := models.SshClientWithPasswd(sshUser, sshPwd, sshIp, models.SshPort)
	if err != nil {
		return fmt.Errorf("clearDelay, create SshClientWithPasswd for ip [%s], error: %w", sshIp, err)
	}
	defer sshClient.Close()

	output, err := models.SshOneCommand(sshClient, clearDelayCmd)
	if err != nil {
		beego.Error(fmt.Sprintf("clearDelay, Execute command [%s] on IP [%s] error: [%s]. We just log but not return this error, because if there is not delay added, there will also be an error.", clearDelayCmd, sshIp, err.Error()))
	}
	beego.Info(fmt.Sprintf("Clearing all network delay on IP [%s], output: [%s].", sshIp, string(output)))
	return nil
}
