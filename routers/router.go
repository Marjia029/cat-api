package routers

import (
	"myproject/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	// Serve static files
	beego.SetStaticPath("/static", "static")

	beego.Router("/", &controllers.MainController{})

	beego.Router("/breed-search", &controllers.BreedSearchController{})

	beego.Router("/voting", &controllers.VotingController{})


}