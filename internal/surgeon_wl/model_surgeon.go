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

type Surgeon struct {

	// Unique identifier of the ambulance
	Id string `json:"id"`

	// Human readable name of the surgeon
	Name string `json:"name"`

	SurgeriesList []SurgeryEntry `json:"surgeriesList,omitempty"`
}