package ipfs

import (
	"context"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"net/http"

	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/constants"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"
)

func IpfsUploadFileByWebApi(apiUrl, filefullpath string) (*string, error) {
	response, err := web.HttpUploadFileByStream(apiUrl, filefullpath)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	fileHash := utils.GetFieldStrFromJson(response, "Hash")
	//logs.GetLogger().Info(carFileHash)

	if fileHash == constants.EMPTY_STRING {
		err := fmt.Errorf("cannot get file hash from response:%s", response)
		//logs.GetLogger().Error(err)
		return nil, err
	}

	return &fileHash, nil
}

func Export2CarFile(apiUrl, fileHash string, carFileFullPath string) error {
	apiUrlFull := utils.UrlJoin(apiUrl, "api/v0/dag/export")
	apiUrlFull = apiUrlFull + "?arg=" + fileHash + "&progress=false"
	carFileContent, err := web.HttpRequest(http.MethodPost, apiUrlFull, "", nil)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	bytesWritten, err := utils.CreateFileWithByteContents(carFileFullPath, carFileContent)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}
	logs.GetLogger().Info(bytesWritten, " bytes have been written to:", carFileFullPath)
	return nil
}

func UploadDirToIpfs(sh *shell.Shell, dirName string) (string, error) {

	dataCid, err := sh.AddDir(dirName)
	if err != nil {
		logs.GetLogger().Error(err)
		return "", err
	}

	return dataCid, nil
}

func DataCidIsDirectory(sh *shell.Shell, dataCid string) (*bool, error) {

	path := utils.UrlJoin("/ipfs/", dataCid)
	stat, err := sh.FilesStat(context.Background(), path)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	isFile := false
	if stat.Type == "directory" {
		isFile = true
	}

	return &isFile, nil
}

func DownloadFromIpfs(sh *shell.Shell, dataCid, outDir string) error {
	return sh.Get(dataCid, outDir)
}
