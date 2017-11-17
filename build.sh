CURDIR="./script"
PRODUCT=$(cd $CURDIR; python -c "import settings; print settings.PRODUCT")
SUBSYS=$(cd $CURDIR; python -c "import settings; print settings.SUBSYS")
MODULE=$(cd $CURDIR; python -c "import settings; print settings.MODULE")
if [ -z "$PRODUCT" ] || [ -z "$SUBSYS" ] || [ -z "$MODULE" ]; then
    echo "Support PRODUCT SUBSYS MODULE PORT in settings.py"
    exit -1
fi
REL_NAME=${PRODUCT}_${SUBSYS}_${MODULE}

mkdir -p output/bin output/conf output/static
cp script/bootstrap.sh script/settings.py output
cp -rf static/* output/static/
chmod +x output/bootstrap.sh
cp -rf conf/* output/conf/

GIT_SHA=`git rev-parse --short HEAD || echo "GitNotFound"`
DATE=`date '+%Y%m%d%H%M%S'`
DATE=`date`
GOVERSION=`go version`
VERSION=${GIT_SHA}
val=$(go version)
ver=$(echo $val | awk -F ' ' '{print $3}' | awk -F '.' '{print $2}')
if [ $ver -gt 4 ]; then
    LINK_OPERATOR="="
else
    LINK_OPERATOR=" "
fi
godep go build -ldflags "-X 'main.COMMITVER${LINK_OPERATOR}${VERSION}' -X 'main.GOVERSION${LINK_OPERATOR}${GOVERSION}' -X 'main.DATE=${DATE}'" -o output/bin/$REL_NAME

