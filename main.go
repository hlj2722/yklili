package main

import (
	"beegostudy/controllers/data"
	"beegostudy/controllers/dml"
	"beegostudy/controllers/platform"
	"beegostudy/util/cron"

	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "root:qweqwe@tcp(60.205.164.3:3306)/beestudy?charset=utf8")
}

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.Data["UserName"] = "HHHHH"
	this.Data["Test"] = TestStatic
	ss := Crontab.Entries()
	for _, s := range ss {
		nextS := s.Job
		//上一次执行时间（本次）
		prevTime := s.Prev.Format("2006-01-02 15:04:05")
		//下一次执行时间
		nextTime := s.Next.Format("2006-01-02 15:04:05")
		fmt.Printf("%T\n", nextS)
		fmt.Println(nextS, prevTime, nextTime)
	}
	this.TplName = "test.html"
}

var TestStatic int
var Crontab = cron.New()

func TestCron() {
	spec := "*/3, *, *, *, *, *"
	//c := cron.New()
	Crontab.AddFunc(spec, func() {
		fmt.Println("-----------哈哈------")
		TestStatic++
	})
	spec2 := "*/8, *, *, *, *, *"
	Crontab.AddFunc(spec2, func() {
		fmt.Println("-----------哈哈22------")
		TestStatic++
	})
	Crontab.Start()
	select {}
}

func main() {
	go TestCron()
	orm.Debug = true                                 //ORM调试模式打开
	beego.BConfig.WebConfig.Session.SessionOn = true //启用Session

	beego.Router("/", &MainController{})
	beego.Router("/login", &platform.LoginController{})
	beego.Router("/register", &platform.RegisterController{})

	//页面控制器
	pages := map[string]beego.ControllerInterface{
		"users": &platform.UsersController{},
		"menus": &platform.MenusController{},
	}

	for name, controller := range pages {
		beego.Router(fmt.Sprintf("/platform/%s", name), controller, "get:Page")
	}

	//数据控制器
	models := map[string]beego.ControllerInterface{
		"menu": &data.MenuController{},
		"user": &data.UserController{},
	}

	for name, controller := range models {
		beego.Router(fmt.Sprintf("/data/%s/:method", name), controller, "*:Get")
	}

	beego.Router("/platform/test", &platform.TestController{})

	//校验用户登录：未登录则重定向到login
	var FilterUser = func(ctx *context.Context) {
		if ctx.Input.Session("User") == nil {
			//如果使用dialog方式会出现弹出窗口被定向到了登录页，这里使用js跳转
			//ctx.Redirect(302, "platform.LoginPage")
			ctx.WriteString(platform.LoginPageScript)
		}
	}

	//dml格式支持直接预览
	beego.AddTemplateExt("dml")
	beego.Router("/:path.dml", &dml.DMLController{})

	beego.InsertFilter("/platform/*", beego.BeforeRouter, FilterUser)
	beego.Run()

}
