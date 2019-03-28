package gokits

import (
    "net/http"
)

//noinspection GoUnusedExportedFunction
func ResponseJson(writer http.ResponseWriter, json string) {
    ResponseContent(writer, json, "application/json", "UTF-8")
}

//noinspection GoUnusedExportedFunction
func ResponseText(writer http.ResponseWriter, text string) {
    ResponseContent(writer, text, "application/json", "UTF-8")
}

//noinspection GoUnusedExportedFunction
func ResponseHtml(writer http.ResponseWriter, html string) {
    ResponseContent(writer, html, "application/json", "UTF-8")
}

func ResponseContent(writer http.ResponseWriter,
    content, contentType, characterEncoding string) {
    writer.Header().Set("Content-Type",
        contentType+"; charset="+characterEncoding)
    _, _ = writer.Write([]byte(content))
}

//noinspection GoUnusedExportedFunction
func ResponseErrorJson(writer http.ResponseWriter, statusCode int, json string) {
    ResponseErrorContent(writer, statusCode, json, "application/json", "UTF-8")
}

//noinspection GoUnusedExportedFunction
func ResponseErrorText(writer http.ResponseWriter, statusCode int, text string) {
    ResponseErrorContent(writer, statusCode, text, "application/json", "UTF-8")
}

//noinspection GoUnusedExportedFunction
func ResponseErrorHtml(writer http.ResponseWriter, statusCode int, html string) {
    ResponseErrorContent(writer, statusCode, html, "application/json", "UTF-8")
}

func ResponseErrorContent(writer http.ResponseWriter, statusCode int,
    errorContent, contentType, characterEncoding string) {
    writer.WriteHeader(statusCode)
    ResponseContent(writer, errorContent, contentType, characterEncoding)
}
