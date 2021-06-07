# Secure Multi Party Computation
This is a bachelor project by Kaare Ø. Kristensen, Mads Hejlesen and Jens S. Nellemann.
## Requires 

Go 1.15

To run experiments the following two modules are used:

github.com/360EntSecGroup-Skylar/excelize/v2 v2.4.0

gonum.org/v1/plot v0.9.0

## Running

To run the projects different experiment, first navigate to the root of this project.
If the experiment is run locally, just uncommented the lines for given experiment in runExperiment() in MPC.go.

### Running distributed
Input the ip adresses of the perticipating computers in the ips variable in makeDistibuted()

```
    var ips = []string{"First ip", "Second ip", "third ip"}
```

Run the command to create config files for all the participants.

```
    go run MPC.go files
```

Set the variable computer nr and amount of multiplications.

```
    numberOfMults := 1000
    computerNr := 1 //If your Ip was index 0 in ips
```

Run the command 
```
    go run MPC.go
```
## Disclaimer
This project is not active after June 18. 2021 and will not be maintained in any way. It was only for educational purpose and should not be used in any other way. 
