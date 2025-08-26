package tools

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

// WeatherRequest represents a weather query request
type WeatherRequest struct {
	City string `json:"city" jsonschema:"description=City name to get weather for"`
	Date string `json:"date" jsonschema:"description=Date in YYYY-MM-DD format (optional)"`
}

// WeatherResponse represents weather information
type WeatherResponse struct {
	City        string `json:"city"`
	Temperature int    `json:"temperature"`
	Condition   string `json:"condition"`
	Date        string `json:"date"`
	Error       string `json:"error,omitempty"`
}

// FlightRequest represents a flight search request
type FlightRequest struct {
	From       string `json:"from" jsonschema:"description=Departure city"`
	To         string `json:"to" jsonschema:"description=Destination city"`
	Date       string `json:"date" jsonschema:"description=Departure date in YYYY-MM-DD format"`
	Passengers int    `json:"passengers" jsonschema:"description=Number of passengers"`
}

// FlightResponse represents flight search results
type FlightResponse struct {
	Flights []Flight `json:"flights"`
	Error   string   `json:"error,omitempty"`
}

type Flight struct {
	Airline   string `json:"airline"`
	FlightNo  string `json:"flight_no"`
	Departure string `json:"departure"`
	Arrival   string `json:"arrival"`
	Price     int    `json:"price"`
	Duration  string `json:"duration"`
}

// HotelRequest represents a hotel search request
type HotelRequest struct {
	City     string `json:"city" jsonschema:"description=City to search hotels in"`
	CheckIn  string `json:"check_in" jsonschema:"description=Check-in date in YYYY-MM-DD format"`
	CheckOut string `json:"check_out" jsonschema:"description=Check-out date in YYYY-MM-DD format"`
	Guests   int    `json:"guests" jsonschema:"description=Number of guests"`
}

// HotelResponse represents hotel search results
type HotelResponse struct {
	Hotels []Hotel `json:"hotels"`
	Error  string  `json:"error,omitempty"`
}

type Hotel struct {
	Name      string   `json:"name"`
	Rating    float64  `json:"rating"`
	Price     int      `json:"price"`
	Location  string   `json:"location"`
	Amenities []string `json:"amenities"`
}

// AttractionRequest represents a tourist attraction search request
type AttractionRequest struct {
	City     string `json:"city" jsonschema:"description=City to search attractions in"`
	Category string `json:"category" jsonschema:"description=Category of attractions (museum, park, landmark, etc.)"`
}

// AttractionResponse represents attraction search results
type AttractionResponse struct {
	Attractions []Attraction `json:"attractions"`
	Error       string       `json:"error,omitempty"`
}

type Attraction struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	OpenHours   string  `json:"open_hours"`
	TicketPrice int     `json:"ticket_price"`
	Category    string  `json:"category"`
}

// NewWeatherTool Mock weather tool implementation
func NewWeatherTool(ctx context.Context) (tool.BaseTool, error) {
	return utils.InferTool("get_weather", "Get weather information for a specific city and date",
		func(ctx context.Context, req *WeatherRequest) (*WeatherResponse, error) {
			if req.City == "" {
				return &WeatherResponse{Error: "City is required"}, nil
			}

			// Mock weather data
			weathers := map[string]WeatherResponse{
				"Beijing":  {City: "Beijing", Temperature: 15, Condition: "Sunny", Date: req.Date},
				"Shanghai": {City: "Shanghai", Temperature: 20, Condition: "Cloudy", Date: req.Date},
				"Tokyo":    {City: "Tokyo", Temperature: 18, Condition: "Rainy", Date: req.Date},
				"Paris":    {City: "Paris", Temperature: 12, Condition: "Overcast", Date: req.Date},
				"New York": {City: "New York", Temperature: 8, Condition: "Snow", Date: req.Date},
			}

			if weather, exists := weathers[req.City]; exists {
				return &weather, nil
			}

			// Random weather for unknown cities
			conditions := []string{"Sunny", "Cloudy", "Rainy", "Overcast"}
			return &WeatherResponse{
				City:        req.City,
				Temperature: rand.Intn(30) + 5, // 5-35°C
				Condition:   conditions[rand.Intn(len(conditions))],
				Date:        req.Date,
			}, nil
		})
}

// NewFlightSearchTool Mock flight search tool implementation
func NewFlightSearchTool(ctx context.Context) (tool.BaseTool, error) {
	return utils.InferTool("search_flights", "Search for flights between cities",
		func(ctx context.Context, req *FlightRequest) (*FlightResponse, error) {
			if req.From == "" || req.To == "" {
				return &FlightResponse{Error: "From and To cities are required"}, nil
			}

			// Mock flight data
			airlines := []string{"Air China", "China Eastern", "China Southern", "United Airlines", "Delta"}

			flights := make([]Flight, 3)
			for i := 0; i < 3; i++ {
				flights[i] = Flight{
					Airline:   airlines[rand.Intn(len(airlines))],
					FlightNo:  fmt.Sprintf("%s%d", airlines[rand.Intn(len(airlines))][:2], rand.Intn(9000)+1000),
					Departure: fmt.Sprintf("%02d:%02d", rand.Intn(24), rand.Intn(60)),
					Arrival:   fmt.Sprintf("%02d:%02d", rand.Intn(24), rand.Intn(60)),
					Price:     rand.Intn(2000) + 500, // $500-2500
					Duration:  fmt.Sprintf("%dh %dm", rand.Intn(12)+1, rand.Intn(60)),
				}
			}

			return &FlightResponse{Flights: flights}, nil
		})
}

// NewHotelSearchTool Mock hotel search tool implementation
func NewHotelSearchTool(ctx context.Context) (tool.BaseTool, error) {
	return utils.InferTool("search_hotels", "Search for hotels in a city",
		func(ctx context.Context, req *HotelRequest) (*HotelResponse, error) {
			if req.City == "" {
				return &HotelResponse{Error: "City is required"}, nil
			}

			// Mock hotel data
			hotelNames := []string{"Grand Hotel", "City Center Inn", "Luxury Resort", "Budget Lodge", "Business Hotel"}
			amenities := [][]string{
				{"WiFi", "Pool", "Gym", "Spa"},
				{"WiFi", "Breakfast", "Parking"},
				{"WiFi", "Pool", "Restaurant", "Bar", "Concierge"},
				{"WiFi", "Breakfast"},
				{"WiFi", "Business Center", "Meeting Rooms"},
			}

			hotels := make([]Hotel, 4)
			for i := 0; i < 4; i++ {
				hotels[i] = Hotel{
					Name:      fmt.Sprintf("%s %s", req.City, hotelNames[rand.Intn(len(hotelNames))]),
					Rating:    float64(rand.Intn(30)+20) / 10.0, // 2.0-5.0
					Price:     rand.Intn(300) + 50,              // $50-350 per night
					Location:  fmt.Sprintf("%s Downtown", req.City),
					Amenities: amenities[rand.Intn(len(amenities))],
				}
			}

			return &HotelResponse{Hotels: hotels}, nil
		})
}

// NewAttractionSearchTool Mock attraction search tool implementation
func NewAttractionSearchTool(ctx context.Context) (tool.BaseTool, error) {
	return utils.InferTool("search_attractions", "Search for tourist attractions in a city",
		func(ctx context.Context, req *AttractionRequest) (*AttractionResponse, error) {
			if req.City == "" {
				return &AttractionResponse{Error: "City is required"}, nil
			}

			// Mock attraction data based on city
			attractionsByCity := map[string][]Attraction{
				"Beijing": {
					{Name: "Forbidden City", Description: "Ancient imperial palace", Rating: 4.8, OpenHours: "8:30-17:00", TicketPrice: 60, Category: "landmark"},
					{Name: "Great Wall", Description: "Historic fortification", Rating: 4.9, OpenHours: "6:00-18:00", TicketPrice: 45, Category: "landmark"},
					{Name: "Temple of Heaven", Description: "Imperial sacrificial altar", Rating: 4.6, OpenHours: "6:00-22:00", TicketPrice: 35, Category: "landmark"},
				},
				"Paris": {
					{Name: "Eiffel Tower", Description: "Iconic iron lattice tower", Rating: 4.7, OpenHours: "9:30-23:45", TicketPrice: 25, Category: "landmark"},
					{Name: "Louvre Museum", Description: "World's largest art museum", Rating: 4.8, OpenHours: "9:00-18:00", TicketPrice: 17, Category: "museum"},
					{Name: "Notre-Dame Cathedral", Description: "Medieval Catholic cathedral", Rating: 4.5, OpenHours: "8:00-18:45", TicketPrice: 0, Category: "landmark"},
				},
				"Tokyo": {
					{Name: "Senso-ji Temple", Description: "Ancient Buddhist temple", Rating: 4.4, OpenHours: "6:00-17:00", TicketPrice: 0, Category: "landmark"},
					{Name: "Tokyo National Museum", Description: "Largest collection of cultural artifacts", Rating: 4.3, OpenHours: "9:30-17:00", TicketPrice: 1000, Category: "museum"},
					{Name: "Ueno Park", Description: "Large public park with museums", Rating: 4.2, OpenHours: "5:00-23:00", TicketPrice: 0, Category: "park"},
				},
			}

			if attractions, exists := attractionsByCity[req.City]; exists {
				// Filter by category if specified
				if req.Category != "" {
					var filtered []Attraction
					for _, attraction := range attractions {
						if attraction.Category == req.Category {
							filtered = append(filtered, attraction)
						}
					}
					return &AttractionResponse{Attractions: filtered}, nil
				}
				return &AttractionResponse{Attractions: attractions}, nil
			}

			// Generate random attractions for unknown cities
			attractionNames := []string{"Central Museum", "City Park", "Historic Square", "Art Gallery", "Cultural Center"}
			categories := []string{"museum", "park", "landmark", "gallery", "cultural"}

			attractions := make([]Attraction, 3)
			for i := 0; i < 3; i++ {
				attractions[i] = Attraction{
					Name:        fmt.Sprintf("%s %s", req.City, attractionNames[rand.Intn(len(attractionNames))]),
					Description: "Popular tourist attraction",
					Rating:      float64(rand.Intn(20)+30) / 10.0, // 3.0-5.0
					OpenHours:   "9:00-17:00",
					TicketPrice: rand.Intn(50),
					Category:    categories[rand.Intn(len(categories))],
				}
			}

			return &AttractionResponse{Attractions: attractions}, nil
		})
}

// GetAllTravelTools returns all travel-related tools
func GetAllTravelTools(ctx context.Context) ([]tool.BaseTool, error) {
	weatherTool, err := NewWeatherTool(ctx)
	if err != nil {
		return nil, err
	}

	flightTool, err := NewFlightSearchTool(ctx)
	if err != nil {
		return nil, err
	}

	hotelTool, err := NewHotelSearchTool(ctx)
	if err != nil {
		return nil, err
	}

	attractionTool, err := NewAttractionSearchTool(ctx)
	if err != nil {
		return nil, err
	}

	return []tool.BaseTool{weatherTool, flightTool, hotelTool, attractionTool}, nil
}
