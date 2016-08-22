# Name
golang的ini配置文件解析器

# Usage
        ...
        import(
            "github.com/mobliyishengyuan/goconf"
            "fmt"
        )
        ...
        func main {
           var config = goconf.GetNewConfig()
           
           var status, err = config.Read("simple.ini")
           
           if !status {
              fmt.Println(err)
              return
           }
           
           var value, status = config.Get("section_1", "key_1")
           
           fmt.Println(value, status)
        }
