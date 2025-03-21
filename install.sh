REPO=hookground-cli
BINARY_NAME="tunnel"

detect_platform() {
  local unameOut="$(uname -s)"
  local arch="$(uname -m)"
  
  case "${unameOut}" in
    Linux*)     os="linux";;
    Darwin*)    os="darwin";;
    CYGWIN*|MINGW*|MSYS*|Windows_NT*) os="windows";;
    *)          os="unknown";;
  esac

  case "${arch}" in
    x86_64)     arch="amd64";;
    arm64|aarch64) arch="arm64";;
    *)          arch="unknown";;
  esac

  if [[ "$os" == "unknown" || "$arch" == "unknown" ]]; then
    echo "Unsupported platform: $unameOut / $arch"
    exit 1
  fi

  echo "${os}-${arch}"
}

TOKEN=""
URL=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --t)
      TOKEN="$2"
      shift 2
      ;;
    --h)
      URL="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

if [[ -z "$TOKEN" || -z "$URL" ]]; then
  echo "Usage: $0 --t <TOKEN> --h <URL>"
  exit 1
fi

PLATFORM=$(detect_platform)
FILENAME="${BINARY_NAME}-${PLATFORM}"

if [[ "$PLATFORM" == *"windows"* ]]; then
  FILENAME="${FILENAME}.exe"
fi

echo "Downloading the binary for ${PLATFORM}..."
DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${FILENAME}"
curl -L -o "${BINARY_NAME}" "${DOWNLOAD_URL}"

if [[ $? -ne 0 ]]; then
  echo "Failed to download the binary from ${DOWNLOAD_URL}"
  exit 1
fi

chmod +x "${BINARY_NAME}"

echo "Running the binary..."
./"${BINARY_NAME}" --t "${TOKEN}" --h="{$URL}"