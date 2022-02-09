# GJG
go-jump-goland tool

Allows you to launch quickly the Goland IDE and open project specified by it name (not path).
Automatically scans ~/go/src/ folder for GIT-projects, scanc installed Goland versions and remembers your choice.

# Usage example: 
~ gjg my-awesome-project // locate ~/go/src/.../my-awesome-project/ folder, open it in IDE
~ gjg -r                 // reinit GJG tool (allows to choose other version of the Goland IDE)

# Installation
~ go install
