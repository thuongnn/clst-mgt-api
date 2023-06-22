package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/wpcodevo/golang-mongodb/services"
	"net/http"
)

type NodeController struct {
	nodeService services.NodeService
}

func NewNodeController(nodeService services.NodeService) NodeController {
	return NodeController{nodeService}
}

func (nc *NodeController) GetRoles(ctx *gin.Context) {
	roles, err := nc.nodeService.GetRoles()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": roles})
}

func (nc *NodeController) GetRolesByNodeId(ctx *gin.Context) {
	nodeId := ctx.Param("nodeId")

	roles, err := nc.nodeService.GetRolesByNodeId(nodeId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": roles})
}

func (nc *NodeController) GetNodes(ctx *gin.Context) {
	currentNodes, err := nc.nodeService.GetNodes()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": currentNodes})
}

func (nc *NodeController) SyncNodes(ctx *gin.Context) {
	err := nc.nodeService.SyncNodes()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}
