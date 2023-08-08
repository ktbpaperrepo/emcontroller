### What is this program for?
In Genetic Algorithm, we can set different **crossover probability** and **mutation probability**. This program is to do some tests to find the best **crossover probability** and **mutation probability**.

### How to use this program?

1. Comment out unnecessary logs to avoid to many logs in the output file.
2. In this folder, run `go build`.
3. Copy (`scp`) the generated binary file `optimizecpmp` to a VM.
4. On the VM execute `nohup ./optimizecpmp > output.log 2>&1 &` to run the program in the background.
5. After this program finishes, the data will be in the file `output.log`. 