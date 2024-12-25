package views

import (
	"context"
	"net/url"

	"codexray/api/views/application"
	"codexray/api/views/applications"
	"codexray/api/views/aws"
	"codexray/api/views/incident"
	"codexray/api/views/inspections"
	"codexray/api/views/integrations"
	"codexray/api/views/logs"
	"codexray/api/views/overview"
	"codexray/api/views/profiling"
	"codexray/api/views/roles"
	"codexray/api/views/tracing"
	"codexray/api/views/users"
	"codexray/clickhouse"
	"codexray/db"
	"codexray/model"
	"codexray/rbac"
	"codexray/timeseries"
)

func Overview(ctx context.Context, ch *clickhouse.Client, w *model.World, view, query string) *overview.Overview {
	return overview.Render(ctx, ch, w, view, query)
}

func Application(w *model.World, app *model.Application) *application.View {
	return application.Render(w, app)
}

func Incident(w *model.World, app *model.Application, i *model.ApplicationIncident) *incident.View {
	return incident.Render(w, app, i)
}

func Profiling(ctx context.Context, ch *clickhouse.Client, app *model.Application, q url.Values, wCtx timeseries.Context) *profiling.View {
	return profiling.Render(ctx, ch, app, q, wCtx)
}

func Tracing(ctx context.Context, ch *clickhouse.Client, app *model.Application, q url.Values, w *model.World) *tracing.View {
	return tracing.Render(ctx, ch, app, q, w)
}

func Logs(ctx context.Context, ch *clickhouse.Client, app *model.Application, q url.Values, w *model.World) *logs.View {
	return logs.Render(ctx, ch, app, q, w)
}

func Inspections(checkConfigs model.CheckConfigs) *inspections.View {
	return inspections.Render(checkConfigs)
}

func Categories(p *db.Project) *applications.CategoriesView {
	return applications.RenderCategories(p)
}

func CustomApplications(p *db.Project) *applications.CustomApplicationsView {
	return applications.RenderCustomApplications(p)
}

func Integrations(p *db.Project) *integrations.View {
	return integrations.Render(p)
}

func AWS(w *model.World) *aws.View {
	return aws.Render(w)
}

func Users(us []*db.User, rs []rbac.Role) *users.View {
	return users.Render(us, rs)
}
func Roles(rs []rbac.Role) *roles.View {
	return roles.Render(rs)
}
