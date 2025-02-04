import socket
import time

# ANSI color codes
RED = "\033[91m"
GREEN = "\033[92m"
BLUE = "\033[94m"
RESET = "\033[0m"  # Reset to default color

def create_grid():
    """Creates an 11x11 grid with roads and an intersection."""
    grid = [["•" for _ in range(11)] for _ in range(11)]

    # Draw roads
    for i in range(11):
        if i == 5:  # Center intersection row
            grid[i] = ["-", "-", "-", "-", "-", "-", "-", "-", "-", "-", "-"]
        else:
            grid[i][5] = "|"  # Vertical road
    grid[5][5]="+"
    return grid

def place_vehicles(grid, state_info):
    """Places vehicles on the grid with color coding."""
    # Map directions to grid positions
    positions = {
        "N": (1, 5),
        "S": (1, 5),
        "E": (5, 1),
        "W": (5, 1)
    }

    vehicle_counter = 1  # To give each vehicle a unique number
    north=1
    south=1
    west=1
    est=1

    # Process normal vehicles
    for direction, vehicles in state_info["vehicles_normal"].items():
        for vehicle in vehicles:
            pos = positions.get(direction, None)
            if pos:
                color = GREEN if state_info["lights_state"][direction] == "GREEN" else RED
                if direction == "N":
                    grid[pos[0] - north][pos[1]] = f"{color}{vehicle['id']}{RESET}"
                    north += 1
                elif direction == "S":
                    grid[pos[0] + south][pos[1]] = f"{color}{vehicle['id']}{RESET}"
                    south += 1
                elif direction == "E":
                    grid[pos[0]][pos[1] + est] = f"{color}{vehicle['id']}{RESET}"
                    est += 1
                elif direction == "W":
                    grid[pos[0]][pos[1] - west] = f"{color}{vehicle['id']}{RESET}"
                    west += 1

    # Process priority vehicles
    for direction, vehicles in state_info["vehicles_priority"].items():
        for vehicle in vehicles:
            pos = positions.get(direction, None)
            if pos:
                grid[pos[0]][pos[1]] = f"{BLUE}{vehicle['id']}{RESET}"  # Always blue for priority
                print('\n the priority vehicle passed')

def print_grid(grid):
    """Prints the 11x11 traffic grid."""
    print("\n".join([" ".join(row) for row in grid]))
    print("\n")

def client_display():
    host = 'localhost'
    port = 8080

    while True:
        try:
            # Create a connection to the server
            with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as client_socket:
                client_socket.settimeout(5)  # Set a timeout for the connection
                client_socket.connect((host, port))
                
                # Receive data from the server
                data = client_socket.recv(4096)
                if data:
                    state_info = eval(data.decode('utf-8'))  # Decode the received data
                    grid = create_grid()  # Generate the grid
                    print('ici')
                    place_vehicles(grid, state_info)  # Place vehicles on the grid
                    print_grid(grid)  # Display the updated intersection
                    
                else:
                    print("Aucune donnée reçue du serveur.")
        except ConnectionRefusedError:
            print("Connexion refusée. Assurez-vous que le serveur est en cours d'exécution.")
            break
        except socket.timeout:
            print("Timeout: Aucune donnée reçue du serveur.")
        except Exception as e:
            print(f"Une erreur s'est produite : {e}")
        
        time.sleep(1)  # Wait for 1 second before requesting the next update

if __name__ == "__main__":
    client_display()
