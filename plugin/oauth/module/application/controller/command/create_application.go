package command

import (
	"encoding/json"
	"fmt"

	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/util"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CreateOauthApplication struct of CreateOauthApplication command
type CreateOauthApplication struct {
	applicationManager ApplicationManager

	flagOwnerID   int
	flagOwnerType string
	flagScope     string
	flagDesc      string
}

// NewCreateOauthApplication return struct of CreateOauthApplication
func NewCreateOauthApplication(applicationManager ApplicationManager) *CreateOauthApplication {
	return &CreateOauthApplication{
		applicationManager: applicationManager,
	}
}

// Use return name of command
func (m *CreateOauthApplication) Use() string {
	return "oauth/application:create"
}

// Short return short description of command
func (m *CreateOauthApplication) Short() string {
	return "Create oauth application"
}

// Example return example of command
func (m *CreateOauthApplication) Example() string {
	return "altair plugin oauth/application:create --owner-id 1 --scope read write --owner-type confidential"
}

// Run run command
func (m *CreateOauthApplication) Run(cmd *cobra.Command, args []string) {
	var ownerID *int
	if m.flagOwnerID != 0 {
		ownerID = util.ValueToPointer(m.flagOwnerID)
	}

	var description *string
	if m.flagDesc != "" {
		description = util.ValueToPointer(m.flagDesc)
	}

	var scope *string
	if m.flagScope != "" {
		scope = util.ValueToPointer(m.flagScope)
	}

	oauthApplicationJSON := entity.OauthApplicationJSON{
		OwnerID:     ownerID,
		OwnerType:   util.ValueToPointer(m.flagOwnerType),
		Description: description,
		Scopes:      scope,
	}

	oauthApplicationJSON, err := m.applicationManager.Create(kontext.Fabricate(), oauthApplicationJSON)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	content, _ := json.Marshal(oauthApplicationJSON)
	fmt.Println("Success creating oauth application:", string(content))
}

// Run run command
func (m *CreateOauthApplication) ModifyFlags(flags *pflag.FlagSet) {
	flags.IntVar(&m.flagOwnerID, "owner-id", 0, "Owner ID, can be nil")
	flags.StringVar(&m.flagOwnerType, "owner-type", "", "Owner Type. Enum: confidential, public")
	flags.StringVar(&m.flagScope, "scope", "", "Scope of the application, separated by space")
	flags.StringVar(&m.flagDesc, "desc", "", "Description of the application")
}
