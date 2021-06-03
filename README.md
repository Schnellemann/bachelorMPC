# Secure Multi Party Computation
This is a bachelor project by Kaare Ø. Kristensen, Mads Hejlesen and Jens S. Nellemann.
## Requires 

Go installed

## Running

To run the projects different experiment, first navigate to the root of this project.
If the experiment is run locally, just uncommented the lines for given experiment in runExperiment() in MPC.go.

### Running distributed
Input the ip adresses of the perticipating computers in the ips variable in makeDistibuted()

```
    var ips = []string{"First ip", "Second ip", "third ip"}
```

Run the command 

```
    go run MPC.go files
```

Now set the variable computer nr and amount of multiplications.

```
    numberOfMults := 1000
    computerNr := 1 //If your Ip was index 0 in ips
```

