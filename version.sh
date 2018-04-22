# Get the version.
version=`git describe --tags --long`
# Write out the package.
cat << EOF > version.go
package hambidgetree

//go:generate bash ./version.sh
var Version = "$version"
EOF