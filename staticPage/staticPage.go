package staticPage

import (
	"os"
	"html/template"
	"github.com/danzar/fcwServer/common"
)


var debug = false
/*
   We load the page/templates from the dir
   then load them into the template system for later refence
 */
func PopulateStaticPages(themeName string) *template.Template  {

	common.LogDebugData("Populating Static Page", debug)
	result := template.New("templates")
	templatesPaths := new([]string)

	basePath := "pages"

	templateFolder, err := os.Open(basePath)
	if err != nil{
		panic(err)
	}
	defer templateFolder.Close()

	//Read all the files from the dir
	templateFolderPathsRaw, err := templateFolder.Readdir(-1)
	if err != nil{
		panic(err)
	}

	//Add the read in files into the a string for later use
	for _,pathInfo := range templateFolderPathsRaw{
		common.LogDebugData("Adding Page:"+ basePath + "/" + pathInfo.Name(), debug)
		*templatesPaths = append(*templatesPaths, basePath + "/" + pathInfo.Name())
	}





	//Load Theme templates
	basePath = "themes/" + themeName
	//get and open the themes folder
	templateFolder, err = os.Open(basePath)
	if err != nil{
		panic(err)
	}
	//Defer close
	defer templateFolder.Close()

	//get all the items in the theme folder
	templateFolderPathsRaw, err = templateFolder.Readdir(-1)
	if err != nil{
		panic(err)
	}
	//add the files to the template list
	for _,pathInfo := range templateFolderPathsRaw{
		common.LogDebugData(basePath + "/" + pathInfo.Name(),debug)
		*templatesPaths = append(*templatesPaths, basePath + "/" + pathInfo.Name())
	}

	result.Delims("!{!","!}!")

	//Parse the files for the results
	_, err = result.ParseFiles(*templatesPaths...)
	if err != nil{
		common.LogData(err.Error())
		panic(err)
	}

	return result
}
