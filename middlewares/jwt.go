package plugin

//
//import (
//    "github.com/tietang/go-utils/errs"
//    "github.com/tietang/props/kvs"
//    "github.com/tietang/zebra"
//    "gopkg.in/square/go-jose.v2/jwt"
//    "net/http"
//    "time"
//)
//
//const (
//    KEY_AUTH_TOKEN  = "auth-token"
//    JWT_SIGNING_KEY = "eyJhbGciOiJIUzI1NiJ9"
//    JWT_ENABLED     = "plugin.jwt.enabled"
//)
//
//type Jwt struct {
//    conf kvs.ConfigSource
//}
//
//func (j *Jwt) Use(s *proxy.HttpProxyServer) {
//    if j.conf.GetBoolDefault(JWT_ENABLED, false) {
//        s.Use(j.handle)
//    }
//}
//
//func (j *Jwt) handle(ctx *proxy.Context) error {
//    t := ctx.Cookie(KEY_AUTH_TOKEN)
//    if t == "" {
//        t = ctx.Request.FormValue(KEY_AUTH_TOKEN)
//    }
//
//    if t == "" {
//        err := errs.NilPointError("jwt tocken is empty")
//        handlerError(ctx, err)
//        return err
//    }
//    tok, err := jwt.ParseSigned(t)
//    if err != nil {
//        handlerError(ctx, err)
//        return err
//    }
//
//    cl := jwt.Claims{}
//    if err := tok.Claims(JWT_SIGNING_KEY, &cl); err != nil {
//        handlerError(ctx, err)
//        return err
//    }
//
//    err = cl.Validate(jwt.Expected{
//        Issuer:  "issuer",
//        Subject: "subject",
//        Time:    time.Now(),
//    })
//    if err != nil {
//        handlerError(ctx, err)
//        return err
//    }
//
//    ctx.Next()
//    return nil
//}
//
//func handlerError(ctx *proxy.Context, err error) {
//    if err != nil {
//        ctx.SetStatusCode(http.StatusUnauthorized)
//        ctx.WriteString("Unauthorized: " + err.Error())
//    }
//}
