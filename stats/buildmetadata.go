package stats

import (
	"html/template"
	"strings"

	"github.com/bgetsug/go-toolbox/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	log "github.com/sirupsen/logrus"
)

var (
	bmd BuildInfo
	env config.Environment
)

type BuildInfo struct {
	ProjectTitle   string
	ProjectName    string
	ProjectVersion string
	VCSBranch      string
	VCSRevision    string
	BuildNumber    string
	BuildTimestamp string
}

func SetBuildMetadata(buildInfo BuildInfo, environment config.Environment) {
	bmd = buildInfo
	bmd.ProjectTitle = strings.Replace(bmd.ProjectTitle, "_", " ", -1)
	env = environment
}

func BuildMetadata(router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ignore Content-Type header
		c.Header("Content-Type", "")

		if bmd.BuildNumber == "" {
			bmd.BuildNumber = "n/a"
		}

		buildMetadata := &gin.H{
			"projectTitle":   bmd.ProjectTitle,
			"projectName":    bmd.ProjectName,
			"projectVersion": bmd.ProjectVersion,
			"vcsBranch":      bmd.VCSBranch,
			"vcsRevision":    bmd.VCSRevision,
			"buildNumber":    bmd.BuildNumber,
			"buildTimestamp": bmd.BuildTimestamp,
		}

		switch c.Request.Header.Get("Accept") {
		case "application/json":
			c.JSON(200, buildMetadata)
		case "application/yaml":
			c.YAML(200, buildMetadata)
		default:
			t := template.New("bmd")
			_, err := t.Parse(`
<!doctype html>
<html>
<head>
    <title>{{ .projectTitle }} :: Build Metadata</title>
</head>
<body>
<code>
	<strong>Project Title:</strong> {{ .projectTitle }}<br/>
    <strong>Project Name:</strong> {{ .projectName }}<br/>
    <strong>Project Version:</strong> {{ .projectVersion }}<br/>
    <strong>VCS Branch:</strong> {{ .vcsBranch }}<br/>
    <strong>VCS Revision:</strong> {{ .vcsRevision }}<br/>
    <strong>Build Number:</strong> {{ .buildNumber }}<br/>
    <strong>Build Timestamp:</strong> {{ .buildTimestamp }}
</code>
</body>
</html>
`)
			if err != nil {
				c.AbortWithError(500, err)
				return
			}

			c.Render(200, render.HTML{
				Template: t,
				Data:     buildMetadata,
			})
		}
	}
}

func LogBanner() {
	log.Infof("Starting %s...", bmd.ProjectTitle)
	log.Infof("Version: %s", bmd.ProjectVersion)
	log.Infof("VCS Branch: %s", bmd.VCSBranch)
	log.Infof("VCS Revision: %s", bmd.VCSRevision)
	log.Infof("Build Number: %s", bmd.BuildNumber)
	log.Infof("Build Timestamp: %s", bmd.BuildTimestamp)
	log.Infof("Environment: %s", env)
}
