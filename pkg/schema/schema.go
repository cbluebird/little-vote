package schema

import (
	"log"

	"github.com/graphql-go/graphql"

	"little-vote/pkg/dao"
	"little-vote/pkg/kafka"
	"little-vote/pkg/ticket"
)

var ticketType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Ticket",
	Description: "Ticket Model",
	Fields: graphql.Fields{
		"ticket": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var queryTicket = graphql.Field{
	Name:        "QueryTicket",
	Description: "Query Ticket",
	Type:        ticketType,
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		ticket.ServerTicket.RLock()
		defer ticket.ServerTicket.RUnlock()
		return map[string]interface{}{
			"ticket": ticket.ServerTicket.TicketId,
		}, nil
	},
}

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "User",
	Description: "User Model",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"count": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var queryUser = graphql.Field{
	Name:        "QueryUser",
	Description: "Query User",
	Type:        userType,
	Args: graphql.FieldConfigArgument{
		"name": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		name, _ := p.Args["name"].(string)
		count, err := dao.GetUserInCache(name)
		if err != nil {
			user, err := dao.GetUserInfo(name)
			if err != nil {
				return nil, err
			}
			err = dao.SetUserInCache(user.Name)
			if err != nil {
				return nil, err
			}
			count = user.Count
		}
		return map[string]interface{}{
			"name":  name,
			"count": count,
		}, nil
	},
}

var voteType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Vote",
	Description: "Vote Model",
	Fields: graphql.Fields{
		"code": &graphql.Field{
			Type: graphql.Int,
		},
		"msg": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var mutationVote = graphql.Field{
	Name:        "Vote",
	Description: "Vote",
	Type:        voteType,
	Args: graphql.FieldConfigArgument{
		"name": &graphql.ArgumentConfig{
			Type: graphql.NewList(graphql.String),
		},
		"ticket": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		names, _ := p.Args["name"].([]interface{})
		ticketStr, _ := p.Args["ticket"].(string)

		ticket.ServerTicket.RLock()
		if ticket.ServerTicket.TicketId != ticketStr {
			ticket.ServerTicket.RUnlock()
			return map[string]interface{}{
				"code": 1,
				"msg":  "fail to check ticket",
			}, nil
		}
		if ticket.ServerTicket.Count >= ticket.MAX {
			ticket.ServerTicket.RUnlock()
			return map[string]interface{}{
				"code": 2,
				"msg":  "over ticket limit",
			}, nil
		}
		ticket.ServerTicket.Count++
		ticket.ServerTicket.RUnlock()

		for _, name := range names {
			nameStr := name.(string)
			_, err := dao.GetUserInCache(nameStr)
			if err != nil {
				err = dao.SetUserInCache(nameStr)
				if err != nil {
					return map[string]interface{}{
						"code": 1,
						"msg":  "fail to set cache",
					}, nil
				}
			}
			err = dao.IncrUserInCache(nameStr)
			if err != nil {
				return map[string]interface{}{
					"code": 1,
					"msg":  "fail to incr",
				}, nil
			}
			err = kafka.Send(nameStr)
			if err != nil {
				log.Println("fail to send message in kafka", err)
			}
		}
		return map[string]interface{}{
			"code": 0,
			"msg":  "success",
		}, nil
	},
}

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name:        "RootQuery",
	Description: "Root Query",
	Fields: graphql.Fields{
		"ticket": &queryTicket,
		"user":   &queryUser,
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name:        "RootMutation",
	Description: "Root Mutation",
	Fields: graphql.Fields{
		"vote": &mutationVote,
	},
})

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
