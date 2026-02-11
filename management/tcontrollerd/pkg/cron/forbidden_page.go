package cron

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"

	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/pkg/constants"
	"chaitin.cn/patronus/safeline-2/management/tcontrollerd/utils"
)

const (
	// SpecUpdatePolicy http://www.quartz-scheduler.org/documentation/quartz-2.3.0/tutorials/tutorial-lesson-06.html
	// Seconds Minutes Hours Day-of-Month Month Day-of-Week Year (optional field)
	// every 30 second (starting from 0s)
	specCheckForbiddenPage = "0/30 * * * * ?"

	forbiddenPagePath = "/etc/nginx/forbidden_pages/default_forbidden_page.html"
)

func checkAndUpdateForbiddenPage() {
	existed, err := utils.FileExist(forbiddenPagePath)
	if err != nil {
		logger.Error(err)
		return
	}
	if !existed {
		err = utils.EnsureWriteFile(forbiddenPagePath, []byte(constants.DefaultForbiddenPage), 0644)
		if err != nil {
			logger.Error(err)
		}
		return
	}

	// Only update if the file is the default one, don't overwrite custom pages
	content, err := ioutil.ReadFile(forbiddenPagePath)
	if err != nil {
		logger.Error(err)
		return
	}

	hash := md5.New()
	hash.Write([]byte(content))
	forbiddenMd5 := hex.EncodeToString(hash.Sum(nil))
	if forbiddenMd5 == constants.DefaultForbiddenPageMd5 {
		return
	}

	// Check if it's a custom page (website_*_forbidden_page.html)
	if strings.HasPrefix(filepath.Base(forbiddenPagePath), "website_") && strings.HasSuffix(filepath.Base(forbiddenPagePath), "_forbidden_page.html") {
		// This is a custom forbidden page, don't overwrite it
		return
	}

	err = utils.EnsureWriteFile(forbiddenPagePath, []byte(constants.DefaultForbiddenPage), 0644)
	if err != nil {
		logger.Error(err)
		return
	}
}
