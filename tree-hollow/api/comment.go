package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"tree-hollow/model"
	"tree-hollow/service"
	"tree-hollow/utils"
)

func getComment(ctx *gin.Context) {
	strId := ctx.Query("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		utils.RespErrorWithDate(ctx, "id invalid!")
		return
	}
	infoType := ctx.DefaultQuery("type", "brief")
	if infoType != "brief" && infoType != "details" {
		utils.RespErrorWithDate(ctx, "type invalid!")
		return
	}

	searched, err := strconv.ParseBool(ctx.DefaultQuery("search_child", "false"))
	if err != nil {
		utils.RespErrorWithDate(ctx, "search_child invalid!")
		return
	}

	if !searched {
		if infoType == "brief" {
			getCommentBrief(ctx, id)
		} else {
			getCommentDetails(ctx, id)
		}
		return
	}

	searchComment(ctx, id)

}

func getCommentBrief(ctx *gin.Context, id int) {
	brief, err := service.GetCommentBrief(id)
	if err != nil {
		utils.RespInternalError(ctx)
		return
	}
	utils.RespSuccessfulWithDate(ctx, brief)
}

func getCommentDetails(ctx *gin.Context, id int) {
	details, err := service.GetCommentDetails(id)
	if err != nil {
		utils.RespInternalError(ctx)
		return
	}
	utils.RespSuccessfulWithDate(ctx, details)
}

func searchComment(ctx *gin.Context, id int) {
	comment, err := service.GetAllChildComment(id)
	if err != nil {
		utils.RespInternalError(ctx)
		return
	}
	utils.RespSuccessfulWithDate(ctx, comment)
}

func updateComment(ctx *gin.Context) {

	name := ctx.GetString("username")

	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		utils.RespErrorWithDate(ctx, "id invalid!")
		return
	}

	checked, err := service.CheckCommentIdMatchName(id, name)
	if err != nil {
		utils.RespInternalError(ctx)
		return
	}

	brief, err := service.GetCommentBrief(id)
	if err != nil {
		utils.RespInternalError(ctx)
		return
	}

	if !checked {
		utils.RespErrorWithDate(ctx, "你不能更新别人的评论")
		return
	}

	content := ctx.PostForm("content")
	isOpen, err := strconv.ParseBool(ctx.PostForm("is_open"))

	if err != nil {
		utils.RespErrorWithDate(ctx, "is_open invalid!")
		return
	}

	comment := model.Comment{
		Id:          id,
		ParentId:    brief.ParentId,
		CommentType: brief.CommentType,
		Content:     content,
		SnitchName:  brief.SnitchName,
		IsOpen:      isOpen,
		CommentTime: brief.CommentTime,
		UpdateTime:  time.Now(),
	}
	err = service.UpdateComment(comment)
	if err != nil {
		utils.RespInternalError(ctx)
		return
	}
	utils.RespSuccessful(ctx)
}

func addComment(ctx *gin.Context) {

	name := ctx.GetString("username")

	parentId, err := strconv.Atoi(ctx.PostForm("parent_id"))
	if err != nil {
		utils.RespErrorWithDate(ctx, "parent_id invalid!")
		return
	}
	commentType, err := strconv.Atoi(ctx.PostForm("comment_type"))
	if err != nil {
		utils.RespErrorWithDate(ctx, "comment_type invalid!")
		return
	}
	content := ctx.PostForm("content")

	isOpen, err := strconv.ParseBool(ctx.PostForm("is_open"))
	if err != nil {
		utils.RespErrorWithDate(ctx, "is_open invalid!")
		return
	}
	comment := model.Comment{
		ParentId:    parentId,
		CommentType: model.CommentType(commentType),
		Content:     content,
		SnitchName:  name,
		IsOpen:      isOpen,
		CommentTime: time.Now(),
		UpdateTime:  time.Now(),
	}
	err = service.AddComment(comment)
	if err != nil {
		utils.RespInternalError(ctx)
		return
	}
	utils.RespSuccessful(ctx)
}

func deleteComment(ctx *gin.Context) {
	name := ctx.GetString("username")

	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		utils.RespErrorWithDate(ctx, "id invalid!")
		return
	}

	checked, err := service.CheckCommentIdMatchName(id, name)
	if err != nil {
		utils.RespInternalError(ctx)
		return
	}

	if !checked {
		utils.RespErrorWithDate(ctx, "你不能删除别人的评论")
		return
	}

	err = service.DeleteComment(id)
	if err != nil {
		utils.RespInternalError(ctx)
		return
	}

	utils.RespSuccessful(ctx)

}
