DOC="
Starting setup for devstack.

It will do the following:
- Print this message and wait for confirmation

- Install the following tools if not already installed (might make changes to .zshrc/.bashrc/...)
    - brew (if not available, will need sudo access to install)
    - kubectl (Kubernetes Cli)
    - werf
    - helmfile
    - devspace (v5.18.5 not latest)
    - python3
    - pbincli (cli for privatebin)
    - go
    - k8s-oidc-helper (https://github.com/micahhausler/k8s-oidc-helper)

- Configure these tools with kubernetes cluster info

- [Needs VPN] [OIDC Login] Configure the tools to use your razorpay email to login to the kubernetes cluster

- [Needs VPN] [Spinnaker Pipeline Trigger] Provision access to the kubernetes cluster for your razorpay email


Make sure you're connected to the VPN.

If you don't have homebrew installed (i.e. running brew --version gives 'command not found'),
use razorpay self-serve app to make yourself admin before running this script again.
"

SHELL_TYPE="$(printf '%s' "$SHELL" | rev | cut -d'/' -f1 | rev)"
SHRC_FILE="${HOME}/.${SHELL_TYPE}rc"

BIN_DIR="${HOME}/.devstack/bin"
BIN_DIR_EXPR="\${HOME}/.devstack/bin"

OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

append_line_to_file() {
    declare line="$1"
    declare file="$2"

    [[ -e $file ]] || touch "$file"

    echo "$line" >> "$file"
}

refresh_shrc_binding() {
    source "$SHRC_FILE"
}

add_cmd_to_shrc() {
    declare cmd="$1"

    grep -qsxF -- "$cmd" "$SHRC_FILE" || append_line_to_file "$cmd" "$SHRC_FILE"
    refresh_shrc_binding
}

check_path_contains() {
    declare dir="$1"

    [[ "$PATH" = *":$dir:"* ]] || [[ "$PATH" = *":$dir" ]] || [[ "$PATH" = "$dir:"* ]]
}

add_dir_to_path() {
    declare pathExpression="$1"
    declare exepectedPathComponent="$2"
    declare pathAppendCmd="export PATH=\"${pathExpression}:\${PATH}\""

    check_path_contains "$exepectedPathComponent" || add_cmd_to_shrc "$pathAppendCmd"
}

install_binary() {
    declare url="$1"
    declare dir="$2"
    declare bin="$3"

    mkdir -p "$dir"
    curl -L "$url" > "$dir/$bin"
    chmod +x "$dir/$bin"
}

install() {
    declare cmdName="$1"
    declare installCmd="${2-}"
    declare versionCmd="${3-}"

    echo "looking for $cmdName"
    declare path="$(which $cmdName)" || true
    if [[ -z "$path" ]]; then
        echo "couldn't find $cmdName. installing..."
        if [[ -z "$installCmd" ]]; then
            # default for installation
            brew install "$cmdName"
        else
            "$installCmd"
        fi
    else
        echo "found $cmdName at $path"
    fi

    if [[ -z "$versionCmd" ]]; then
        # default for version check
        "$cmdName" --version
    else
        "$versionCmd"
    fi
}

abort() {
    declare message="$1"

    echo "$message"
    exit 1
}

read_email() {
    declare target="$1"

    read -p "Enter your (razorpay) email address:" "$target"
    is_rzp_email ${!target} || abort "Not a valid razorpay email address"
}

confirm() {
    declare prompt="$1"
    
    read -p "${prompt}Press enter to continue. Press any other key to stop." -n 1

    [[ -z $REPLY ]]
}

spinnaker_webhook() {
    declare spinnaker="$1"
    declare webhook="$2"
    declare parameters="$3"

    curl -X POST "https://$spinnaker/webhooks/webhook/$webhook" \
        -H "content-type: application/json" \
        -d "{\"parameters\":$parameters}"
}

is_rzp_email() {
    declare input="$1"

    [[ "$input" =~ ^[a-zA-Z0-9.!\#$%\&\'*+/=?^_\`{|}~-]+@razorpay\.com$ ]]
}

oidc_exists() {
    declare email="$1"

    declare template="{{\$res := 0}}{{if .users}}{{range .users}}{{if eq .name \"$email\" }}{{\$res = 1}}{{end}}{{end}}{{end}}{{\$res}}"

    [[ $(kubectl config view -o=go-template --template="$template") == 1 ]]
}

install_devspace() {
    declare tag="${OS}-${ARCH}"
    declare version="v5.18.5"
    declare url="https://github.com/loft-sh/devspace/releases/download/$version/devspace-$tag"

    install_binary "$url" "${BIN_DIR}" "devspace"
}

install_brew() {
    /bin/bash -e -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    add_cmd_to_shrc 'eval "$(/opt/homebrew/bin/brew shellenv)"'
}

version_brew() {
    brew config
}

version_kubectl() {
    kubectl version --client --output yaml
}

install_werf() {
    declare tag="${OS}-${ARCH}"
    declare version="1.2.174"
    declare url="https://tuf.werf.io/targets/releases/$version/$tag/bin/werf"

    install_binary "$url" "${BIN_DIR}" "werf"
}

version_werf() {
    echo "werf version: $(werf version)"
}

install_pbincli() {
    pip3 install pbincli
}

version_pbincli() {
    pbincli --help
}

version_go() {
    go version
}

configure_helmfile_for_werf() {
    add_cmd_to_shrc "export WERF_HELM3_MODE=1"
    add_cmd_to_shrc "alias helmfile='helmfile --enable-live-output -b werf'"
}

install_oidc_helper() {
    go install github.com/micahhausler/k8s-oidc-helper@latest
}

cluster_config() {
    declare name="$1"
    declare server="$2"
    declare user="$3"

    kubectl config set-cluster "$name" --server=$server --insecure-skip-tls-verify=true
    kubectl config set-context "$name" --cluster="$name" --user="$user"
    kubectl config use-context "$name"

    echo "kubectl config current-context : $(kubectl config current-context)"
}

oidc_config() {
    declare email="$1"
    declare pasteUrl="$2"
    declare pasteFile="$3"

    oidc_exists "$email" && return 0
    pbincli get "$pasteUrl"
    k8s-oidc-helper -c "./$pasteFile" --write
    rm "./$pasteFile"
}

setup_go_bin() {
    declare goBinDir="$(go env GOPATH)/bin"
    declare goBinDirExpr="\$(go env GOPATH)/bin"

    add_dir_to_path "$goBinDir" "$goBinDirExpr"
}

setup_python_bin() {
    declare pythonBinDir="$(python3 -m site --user-base)/bin"
    declare pythonBinDirExpr="\$(python3 -m site --user-base)/bin"

    add_dir_to_path "$pythonBinDir" "$pythonBinDirExpr"
}

setup_tools() {
    install "brew" "install_brew" "version_brew"
    install "kubectl" "" "version_kubectl"
    install "helmfile"

    add_dir_to_path "${BIN_DIR_EXPR}" "${BIN_DIR}"
    install "werf" "install_werf" "version_werf"
    install "devspace" "install_devspace"

    configure_helmfile_for_werf

    install "python3"
    setup_python_bin

    install "go" "" "version_go"
    setup_go_bin

    install "pbincli" "install_pbincli" "version_pbincli"
    install "k8s-oidc-helper" "install_oidc_helper"
}
