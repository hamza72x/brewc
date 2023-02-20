#/bin/bash

# This script will install the latest version of the `brewc` into /usr/local/bin

# Check if the user has curl, jq, cut installed
green='\033[0;32m'
end='\033[0m'

echo -e "${green}Checking if curl, jq, cut are installed...${end}"

if ! command -v curl &> /dev/null
then
    echo "curl could not be found"
    exit
fi

if ! command -v jq &> /dev/null
then
    echo "jq could not be found"
    exit
fi

if ! command -v cut &> /dev/null
then
    echo "cut could not be found"
    exit
fi

# lowercase uname output
arch=$(uname -m | tr '[:upper:]' '[:lower:]')
goos=$(uname -s | tr '[:upper:]' '[:lower:]')

if [ "$arch" = "x86_64" ]; then
    arch="amd64"
fi

latest_tag=$(curl -sL https://api.github.com/repos/hamza72x/brewc/tags | jq -r '.[0].name' | cut -d 'v' -f 2)
url=https://github.com/hamza72x/brewc/releases/download/v$latest_tag/brewc_$latest_tag\_$goos\_$arch.tar.gz

echo -e "${green}Detected OS: $goos${end}"
echo -e "${green}Detected Arch: $arch${end}"

echo ""

echo -e "${green}Latest version: $latest_tag${end}"
echo -e "${green}Download URL: $url${end}"
echo -e "${green}Installing brewc into /usr/local/bin${end}"


echo ""

# now ask if everything looks good and want to continue
read -p "Looks good? Continue? [y/N] " -n 1 -r

# quit if the user doesn't want to continue
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo
    echo "Aborting..."
    exit 1
fi

# download the latest version
curl -L $url | tar -xz -C /usr/local/bin

# print usage
echo ""
echo -e "${green}Verifying installation...${end}"
echo ""

if command -v brewc &> /dev/null
then
    brewc --help
else
    echo "brewc could not be found, make sure /usr/local/bin is in your PATH, or try /usr/local/bin/brewc --help"
    exit
fi
