set -eu

CLASP=`npm bin`/clasp

DEPLOY_NAME=$1
DEPLOY_ID=`${CLASP} deployments | grep ${DEPLOY_NAME} | cut -d " " -f 2`

if [[ "${DEPLOY_ID}" = "" ]]; then
	DEPLOY_ID=`${CLASP} deploy -d ${DEPLOY_NAME} | grep -v Created | cut -d " " -f 2`
else
	${CLASP} deploy -i ${DEPLOY_ID}
fi

echo https://script.google.com/macros/s/${DEPLOY_ID}/exec
