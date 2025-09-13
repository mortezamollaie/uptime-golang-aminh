# API Documentation

This project includes comprehensive Swagger/OpenAPI documentation for all endpoints.

## Accessing Swagger UI

Once the server is running, you can access the interactive API documentation at:

```
http://localhost:3000/swagger/index.html
```

## Available Endpoints

### Health Check
- `GET /health` - Check API health status

### Node Management
- `GET /api/nodes` - Get all monitoring nodes
- `POST /api/nodes` - Create a new monitoring node
- `GET /api/nodes/{id}` - Get specific node
- `PUT /api/nodes/{id}` - Update node
- `DELETE /api/nodes/{id}` - Delete node

### Reports & Analytics
- `GET /api/report/get` - Get node monitoring report
- `GET /api/report/get-smart-query` - Smart query report
- `POST /api/report/bulk-url/get` - Bulk URL reports
- `GET /api/report/all-from-history` - Complete history
- `GET /api/report/last` - Latest reports

### Node Logs
- `GET /api/node-logs/{id}` - Get node logs
- `GET /api/uptime/{id}` - Get uptime statistics

## Authentication

Most endpoints require API key authentication. Include your API key in the request header:

```
X-API-Key: your_api_key_here
```

Or as Authorization header:

```
Authorization: your_api_key_here
```

## Regenerating Documentation

To regenerate the Swagger documentation after making changes:

```bash
swag init -g cmd/main.go
```

This will update the files in the `docs/` directory:
- `docs.go` - Go package with embedded documentation
- `swagger.json` - OpenAPI JSON specification
- `swagger.yaml` - OpenAPI YAML specification

## Adding New Endpoints

To document new endpoints, add Swagger annotations above your controller functions:

```go
// @Summary Short description
// @Description Detailed description
// @Tags endpoint-group
// @Accept json
// @Produce json
// @Param name query string false "Parameter description"
// @Success 200 {object} ResponseType "Success description"
// @Failure 400 {object} map[string]string "Error description"
// @Security ApiKeyAuth
// @Router /api/endpoint [get]
func MyEndpoint(c *fiber.Ctx) error {
    // Your code here
}
```

## Response Examples

All API responses follow a consistent structure:

### Success Response
```json
{
  "code": 200,
  "msg": "Success",
  "success": true,
  "data": { ... }
}
```

### Error Response
```json
{
  "code": 400,
  "msg": "Error message",
  "success": false,
  "data": null
}
```
