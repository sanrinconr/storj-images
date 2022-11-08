#!/bin/bash
mkdir -p coverage

# Var definitions
coverProfileRaw="coverage/profile.out.raw"
coverProfileFiltered="coverage/profile.out.filtered"
result="coverage/result.raw"

# Execute tests and generate a profile
go test ./... -coverprofile $coverProfileRaw -covermode count

# Delete packages ignored
cat $coverProfileRaw | grep -v /infrastructure/ | grep -v /controller/add_image.go >  $coverProfileFiltered

# Process results, generate results
go tool cover -func $coverProfileFiltered > $result
cat $result

# Extract manually % of coverage
totalCoverage=`cat $result | grep total | grep -Eo '[0-9]+\.[0-9]+'`

# Validate if coverage pass
expectedCoverage=90
if [ 1 -eq "$(echo "$totalCoverage < $expectedCoverage" |bc )" ]
then
    echo "coverage was $totalCoverage% and is needed $expectedCoverage% "
    go tool cover -html=$coverProfileFiltered
    exit 1
fi
