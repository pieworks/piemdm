# API Reference

The PieMDM Backend API serves as the core of the platform, providing data management capabilities.

## Base URL

By default, the API is accessible at:

```
http://127.0.0.1:8787/swagger/index.html
```

## Technology Stack

The backend is built using:
- **Language**: Go
- **Framework**: Gin
- **Database**: GORM (MySQL)
- **Cache**: Redis

## Endpoints

> [!NOTE]
> Detailed API documentation (Swagger/OpenAPI) is typically generated from the code. Please check the backend repository for `swagger` or `docs` folders.

Common modules usually include:
- **Auth**: User authentication and token management.
- **User**: User profile and management.
- **System**: System configuration and status.
