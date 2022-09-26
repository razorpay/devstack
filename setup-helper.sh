set -euo pipefail

DOC="
Starting setup for devstack
This expects bash to be available
Will do the following
- Print this message and wait for confirmation
- Install the following tools if not already installed.
  Note: this adds path related commands and aliases to your profile files (.zshrc, .bashrc, ...)
    - brew
    - kubectl (Kubernetes Cli)
    - werf
    - gh (Github Cli)
    - helmfile (https://github.com/razorpay/helmfile)
    - devspace (v5.18.5 not latest)
    - python3
    - pbincli (cli for privatebin)
    - go
    - k8s-oidc-helper (https://github.com/micahhausler/k8s-oidc-helper)
- [Needs VPN] [OIDC Login] Configure the tools to use your razorpay email to login to the kubernetes cluster
- Configure the tools with kubernetes cluster info
- [Needs VPN] [Spinnaker Pipeline Trigger] Provision access to the kubernetes cluster for your razorpay email

Make sure you're connected to the VPN before continuing.
"

SHELL_TYPE="$(printf '%s' "$SHELL" | rev | cut -d'/' -f1 | rev)"
SHRC_FILE="${HOME}/.${SHELL_TYPE}rc"

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

brew_install_from_url() {
    declare formula="$1"
    declare url="$2"

    # command || true allows us to suppress errors in command
    brew unpin $formula || true
    brew uninstall $formula || true

    formulaPath="$(find $(brew --repository)/Library -name $formula.rb)"

    curl $url > $formulaPath && brew reinstall $formula

    pwdBefore=$(pwd) && cd "$(dirname $formulaPath)" && git checkout . && cd "$pwdBefore"
}

install_devspace() {
    declare sourceCommit="eefcf5566171216e59c89a8c9cf88d38b97f4c74"
    declare sourceUrl="https://raw.githubusercontent.com/Homebrew/homebrew-core/$sourceCommit/Formula/devspace.rb"
    brew_install_from_url "devspace" $sourceUrl
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
            brew install "$cmdName"
        else
            "$installCmd"
        fi
    else
        echo "found $cmdName at $path"
    fi

    if [[ -z "$versionCmd" ]]; then
        "$cmdName" --version
    else
        "$versionCmd"
    fi
}

install_brew() {
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
}

version_brew() {
    brew config
}

version_kubectl() {
    kubectl version --client --output yaml
}

install_helmfile() {
    declare repo="razorpay/helmfile"
    declare tag="v0.144.0-razorpay"
    declare file="helmfile_$(uname | tr '[:upper:]' '[:lower:]')_$(uname -m)"

    [[ -f "${HOME}/bin/$file" ]] || gh release download -R $repo $tag -p "$file" -D "${HOME}/bin"
    chmod +x "${HOME}/bin/$file"
    ln -s "${HOME}/bin/$file" "${HOME}/bin/helmfile"
}

install_werf() {
    declare installer="/tmp/werf-install.sh"
    [[ -f "$installer" ]] || curl -sSL https://werf.io/install.sh -o "$installer"
    chmod +x "$installer"
    "$installer" --version 1.2 --channel stable --no-interactive
    source "$($HOME/bin/trdl use werf 1.2 stable)"
    rm "$installer"
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

install_go() {
    brew install go
}

version_go() {
    go version
}

configure_helmfile_for_werf() {
    add_cmd_to_shrc "export WERF_HELM3_MODE=1"
    add_cmd_to_shrc "alias helmfile='helmfile -b werf --runner-skip-prefix --runner-log-level=info'"
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

welcome() {
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

oidc_config() {
    declare email="$1"
    declare pasteUrl="$2"
    declare pasteFile="$3"

    oidc_exists "$email" && return 0
    pbincli get "$pasteUrl"
    k8s-oidc-helper -c "./$pasteFile" --write
    rm "./$pasteFile"
}

setup_tools() {
    add_dir_to_path "\${HOME}/bin" "${HOME}/bin"
    install "brew" "install_brew" "version_brew"
    install "kubectl" "" "version_kubectl"
    install "werf" "install_werf" "version_werf"
    install "gh"
    install "helmfile" "install_helmfile"
    install "devspace" "install_devspace"
    install "python3"
    install "pbincli" "install_pbincli" "version_pbincli"
    install "go" "install_go" "version_go"
    add_dir_to_path "\$(go env GOPATH)/bin" "$(go env GOPATH)/bin"
    install "k8s-oidc-helper" "install_oidc_helper"
    configure_helmfile_for_werf
    refresh_shrc_binding
}

oidc_exists() {
    declare email="$1"

    [[ "$email" == $(kubectl config view -o jsonpath="{.users[?(@.name == \"$email\")].name}") ]]
}


is_email() {
    declare input="$1"

    [[ "$input" =~ ^[a-zA-Z0-9.!\#$%\&\'*+/=?^_\`{|}~-]+@razorpay\.com$ ]]
}

abort() {
    echo "$1"
    exit 1
}
