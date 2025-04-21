package http

import "github.com/gofiber/fiber/v2"

type Router struct {
	group fiber.Router
}

func (r *Router) Group(path string, groupFunc func(*Router)) *Router {
	groupRouter := &Router{
		group: r.group.Group(path),
	}

	groupFunc(groupRouter)

	return groupRouter
}

func (r *Router) Get(path string, handler fiber.Handler) {
	r.group.Get(path, handler)
}

func (r *Router) Post(path string, handler fiber.Handler) {
	r.group.Post(path, handler)
}

func (r *Router) Put(path string, handler fiber.Handler) {
	r.group.Put(path, handler)
}

func (r *Router) Delete(path string, handler fiber.Handler) {
	r.group.Delete(path, handler)
}

func (r *Router) Patch(path string, handler fiber.Handler) {
	r.group.Patch(path, handler)
}

func (r *Router) Options(path string, handler fiber.Handler) {
	r.group.Options(path, handler)
}

func (r *Router) Head(path string, handler fiber.Handler) {
	r.group.Head(path, handler)
}
