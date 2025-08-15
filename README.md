# Order Packs Calculator

A Go application that calculates optimal pack combinations for customer orders, following specific business 
rules for minimizing items and pack counts.

## Business Rules

The application follows these three rules in order of precedence:

1. **Rule 1**: Only whole packs can be sent. Packs cannot be broken open.
2. **Rule 2**: Send the least amount of items to fulfill the order (takes precedence).
3. **Rule 3**: Send as few packs as possible to fulfill each order.

## Examples

| Items Ordered | Correct Pack Combination | Total Items | Total Packs |
|---------------|-------------------------|-------------|-------------|
| 1 | 1 × 250 | 250 | 1 |
| 250 | 1 × 250 | 250 | 1 |
| 251 | 1 × 500 | 500 | 1 |
| 501 | 1 × 500 + 1 × 250 | 750 | 2 |
| 12001 | 2 × 5000 + 1 × 2000 + 1 × 250 | 12250 | 4 |

## API Endpoints

### Calculate Packs
- **GET** `/calculate?items=501`
- **POST** `/calculate` with body `{"items": 501}`

Response:
```json
{
  "total_items": 750,
  "packs_used": {
    "250": 1,
    "500": 1
  }
}
```

### Manage Pack Sizes
- **GET** `/pack-sizes` - Get current pack sizes
- **POST** `/pack-sizes` - Update pack sizes

Request body:
```json
{
  "pack_sizes": [250, 500, 1000, 2000, 5000]
}
```

## Running the Application

### Local Development

If you choose go through with local installation it's **_your responsibility to make the golang works
on your machine._**.

**I strong suggest to skip to docker way.**

1. **Prerequisites**
   - Go 1.21 or later
   - Git

2. **Installation**
   ```bash
   git clone <repository-url>
   cd packs-mono-repo
   go mod download
   ```

3. **Run the application**
   ```bash
   go run main.go
   ```

4. **Access the application**
   - Web UI: http://localhost:8080
   - API: http://localhost:8080/calculate

### Docker Deployment (preferred)

1. **Build and run with Docker Compose**
   ```bash
   docker-compose up --build
   ```

2. **Or build and run manually**
   ```bash
   docker build -t pack-calculator .
   docker run -p 8080:8080 pack-calculator
   ```

## Project Structure

```
packs-mono-repo/
├── main.go                 # Main application entry point
├── go.mod                  # Go module dependencies
├── go.sum                  # Dependency checksums
├── packs.json              # Default pack sizes configuration
├── public/
│   └── index.html          # Web UI
├── internal/
│   └── packages/           # Core pack calculation logic
├── Dockerfile              # Docker configuration
├── docker-compose.yml      # Docker Compose configuration
└── README.md               # This file
```

## Technology Stack

- **Backend**: Go with Fiber web framework
- **Frontend**: Vanilla JavaScript with modern CSS
- **Containerization**: Docker with multi-stage builds
- **Deployment**: Docker Compose for easy orchestration

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o pack-calculator main.go
```

### Code Structure
- `main.go`: HTTP server setup and routing
- `internal/packages/`: Core algorithm implementation
- `public/index.html`: Modern, responsive web interface

## Algorithm

The application uses an exhaustive search algorithm to find the optimal pack combination:

1. **Single Pack Check**: First tries to find a single pack that can fulfill the order
2. **Two Pack Combinations**: Searches combinations of two different pack sizes
3. **Three Pack Combinations**: Extends to three pack combinations if needed
4. **Four and Five Pack Combinations**: For complex cases, searches up to five pack combinations

The algorithm prioritizes minimizing total items (Rule 2) and then minimizes the number of packs (Rule 3).

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
