package http

// import (
// 	"encoding/json"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/kodefluence/altair/provider/plugin/oauth/entity"
// 	"github.com/kodefluence/altair/provider/plugin/oauth/eobject"
// 	"github.com/kodefluence/altair/provider/plugin/oauth/interfaces"
// 	"github.com/rs/zerolog"
// 	"github.com/rs/zerolog/log"
// )

// // CreateController control flow of oauth application creation
// type CreateController struct {
// 	applicationManager interfaces.ApplicationManager
// }

// // NewCreate return struct of CreateController
// func NewCreate(applicationManager interfaces.ApplicationManager) *CreateController {
// 	return &CreateController{
// 		applicationManager: applicationManager,
// 	}
// }

// // Method POST
// func (cr *CreateController) Method() string {
// 	return "POST"
// }

// // Path /oauth/applications
// func (cr *CreateController) Path() string {
// 	return "/oauth/applications"
// }

// // Control creation of oauth application
// func (cr *CreateController) Control(c *gin.Context) {
// 	var oauthApplicationJSON entity.OauthApplicationJSON

// 	rawData, err := c.GetRawData()
// 	if err != nil {
// 		log.Error().
// 			Err(err).
// 			Stack().
// 			Interface("request_id", c.Value("request_id")).
// 			Array("tags", zerolog.Arr().Str("controller").Str("application").Str("create").Str("get_raw_data")).
// 			Msg("Cannot get raw data")
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"errors": eobject.Wrap(eobject.BadRequestError("request body")),
// 		})
// 		return
// 	}

// 	err = json.Unmarshal(rawData, &oauthApplicationJSON)
// 	if err != nil {
// 		log.Error().
// 			Err(err).
// 			Stack().
// 			Interface("request_id", c.Value("request_id")).
// 			Array("tags", zerolog.Arr().Str("controller").Str("application").Str("update").Str("unmarshal")).
// 			Msg("Cannot unmarshal json")
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"errors": eobject.Wrap(eobject.BadRequestError("request body")),
// 		})
// 		return
// 	}

// 	result, entityError := cr.applicationManager.Create(c, oauthApplicationJSON)
// 	if entityError != nil {
// 		c.JSON(entityError.HttpStatus, gin.H{
// 			"errors": entityError.Errors,
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"data": result,
// 	})
// }
