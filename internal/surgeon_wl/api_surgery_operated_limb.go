/*
 * Waiting List Api
 *
 * Surgeons List for Web-In-Cloud system
 *
 * API version: 1.0.0
 * Contact: tomas.nemec1999@gmail.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package surgeon_wl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SurgeryOperatedLimbAPI interface {

   // internal registration of api routes
   addRoutes(routerGroup *gin.RouterGroup)

    // GetOperatedLimbList - Provides the list of operated limbs associated with surgeries
   GetOperatedLimbList(ctx *gin.Context)

 }

// partial implementation of SurgeryOperatedLimbAPI - all functions must be implemented in add on files
type implSurgeryOperatedLimbAPI struct {

}

func (this *implSurgeryOperatedLimbAPI) GetOperatedLimbList(ctx *gin.Context) {
  // Implementation or dummy function
}

func (this *implSurgeryOperatedLimbAPI) addRoutes(routerGroup *gin.RouterGroup) {
  routerGroup.Handle( http.MethodGet, "/surgeries-list/operatedLimbList", this.GetOperatedLimbList)
}

// Copy following section to separate file, uncomment, and implement accordingly
// // GetOperatedLimbList - Provides the list of operated limbs associated with surgeries
// func (this *implSurgeryOperatedLimbAPI) GetOperatedLimbList(ctx *gin.Context) {
//  	ctx.AbortWithStatus(http.StatusNotImplemented)
// }
//

