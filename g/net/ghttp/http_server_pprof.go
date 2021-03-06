// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// pprof封装.

package ghttp

import (
    "strings"
    runpprof "runtime/pprof"
    netpprof "net/http/pprof"
    "gitee.com/johng/gf/g/frame/gins"
)

// 用于pprof的对象
type utilpprof struct {}

func (p *utilpprof) Index(r *Request) {
    profiles := runpprof.Profiles()
    action   := r.Get("action")
    data     := map[string]interface{}{
        "uri"      : strings.TrimRight(r.URL.Path, "/") + "/",
        "profiles" : profiles,
    }
    if len(action) == 0 {
        view      := gins.View()
        buffer, _ := view.ParseContent(`
            <html>
            <head>
                <title>gf ghttp pprof</title>
            </head>
            {{$uri := .uri}}
            <body>
                profiles:<br>
                <table>
                    {{range .profiles}}<tr><td align=right>{{.Count}}<td><a href="{{$uri}}{{.Name}}?debug=1">{{.Name}}</a>{{end}}
                </table>
                <br><a href="{{$uri}}goroutine?debug=2">full goroutine stack dump</a><br>
            </body>
            </html>
            `, data)
        r.Response.Write(buffer)
        return
    }
    for _, p := range profiles {
        if p.Name() == action {
            p.WriteTo(r.Response.Writer, r.GetRequestInt("debug"))
            break
        }
    }
}

func (p *utilpprof) Cmdline(r *Request) {
    netpprof.Cmdline(r.Response.Writer, &r.Request)
}

func (p *utilpprof) Profile(r *Request) {
    netpprof.Profile(r.Response.Writer, &r.Request)
}

func (p *utilpprof) Symbol(r *Request) {
    netpprof.Symbol(r.Response.Writer, &r.Request)
}

func (p *utilpprof) Trace(r *Request) {
    netpprof.Trace(r.Response.Writer, &r.Request)
}

// 开启pprof支持
func (s *Server) EnablePprof(pattern...string) {
    p := "/debug/pprof"
    if len(pattern) > 0 {
        p = pattern[0]
    }
    up := &utilpprof{}
    _, _, uri, _ := s.parsePattern(p)
    uri = strings.TrimRight(uri, "/")
    s.BindHandler(uri + "/*action", up.Index)
    s.BindHandler(uri + "/cmdline", up.Cmdline)
    s.BindHandler(uri + "/profile", up.Profile)
    s.BindHandler(uri + "/symbol",  up.Symbol)
    s.BindHandler(uri + "/trace",   up.Trace)
}