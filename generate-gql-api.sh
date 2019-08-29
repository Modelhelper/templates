#!/usr/bin/env bash

#mkdir 
#dotnet new webapi 

echo -n "Set root namespace: "
read rns

# printf "\nApi: %s.Api" $rns
# printf "\nData: %s.Data" $rns
# printf "\nCore: %s.Core\nThis assembly is used for infrastructure (interfaces, models, viewmodels)" $rns
# printf "\nTest: %s.Tests\n\tUses XUnit to run tests\n\n" $rns

echo -n "Continue?: [y/n]: "
read cont

# create structure
mkdir $rns
mkdir $rns/docs
mkdir $rns/src
mkdir $rns/src/$rns.Api
mkdir $rns/src/$rns.Api/Extensions
mkdir $rns/src/$rns.Data
mkdir $rns/src/$rns.Core
mkdir $rns/src/$rns.Core/Extensions
mkdir $rns/src/$rns.Core/Models
mkdir $rns/src/$rns.Core/ViewModels
mkdir $rns/test
mkdir $rns/test/$rns.Tests
mkdir $rns/test/$rns.Moc

# root
cd $rns
dotnet new sln 

#root/src/core
cd src/$rns.Core
dotnet new classlib --no-restore

#root/src/data
cd ../$rns.Data
dotnet new classlib --no-restore

#root/src/api
cd ../$rns.Api
dotnet new webapi --no-restore

#root/test/tests
cd ../../test/$rns.Tests
dotnet new xunit --no-restore

#root/test/tests
cd ../$rns.Moc
dotnet new classlib --no-restore



cd ../../
dotnet sln $rns.sln add src/$rns.Api/$rns.Api.csproj src/$rns.Data/$rns.Data.csproj src/$rns.Core/$rns.Core.csproj 
dotnet sln $rns.sln add test/$rns.Tests/$rns.Tests.csproj test/$rns.Moc/$rns.Moc.csproj
#dotnet sln $rns.sln add **/*.csproj

dotnet add src/$rns.Data/$rns.Data.csproj reference src/$rns.Core/$rns.Core.csproj
dotnet add src/$rns.Api/$rns.Api.csproj reference src/$rns.Core/$rns.Core.csproj src/$rns.Data/$rns.Data.csproj

dotnet add test/$rns.Tests/$rns.Tests.csproj reference src/$rns.Core/$rns.Core.csproj test/$rns.Moc/$rns.Moc.csproj
dotnet add test/$rns.Moc/$rns.Moc.csproj reference src/$rns.Core/$rns.Core.csproj

dotnet restore

## add external dependencies
dotnet add src/$rns.Data/$rns.Data.csproj package dapper
dotnet add src/$rns.Api/$rns.Api.csproj package dapper

dotnet build

# cd $rns

# touch .gitignore

# dotnet new sln 

# cd src

# cd "$rns.Api"
# dotnet new webapi --no-restore

# cd ".."

# if ["$cont" == "y"] 
# then    
#     echo " GO "

# fi C:\dev\mh_demo\generate-api.sh