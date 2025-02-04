import socket
import signal
import os
import multiprocessing
import random
import time
import sys

queue_lock = multiprocessing.Lock()

message_queue_priority = {
        "N": multiprocessing.Queue(),
        "S": multiprocessing.Queue(),
        "E": multiprocessing.Queue(),
        "W": multiprocessing.Queue()
    }

def lights(lights_state):
    while True:
        if lights_state["N"]=="GREEN":
            lights_state["N"]="RED"
            lights_state["E"]="GREEN"
            lights_state["S"]="RED"
            lights_state["W"]="GREEN"
        else:
            lights_state["N"]="GREEN"
            lights_state["E"]="RED"
            lights_state["S"]="GREEN"
            lights_state["W"]="RED"
        time.sleep(5)
        print("changing lights")

def route_aleatoire(number_destination):
    first=random.randint(1,4)
    second=random.randint(1,4)
    while first==second:
        second=random.randint(1,4)
    return number_destination[first],number_destination[second]

def generate_id(id_counter):
    with id_counter.get_lock():
        id_counter.value+=1
        return id_counter.value

def normal_traffic_gen(message_queue_normal,number_destination,id_counter):
    while id_counter.value<10:
        entry,dest=route_aleatoire(number_destination)
        id=generate_id(id_counter)
        vehicule={
            "id":id,
            "entry":entry,
            "dest":dest,
            "priority":False
        }
        message_queue_normal[entry].put(vehicule)
        time.sleep(1)

def priority_traffic_gen(message_queue_priority,number_destination,id_counter):
    entry,dest=route_aleatoire(number_destination)
    id=generate_id(id_counter)
    vehicule={
        "id":id,
        "entry":entry,
        "dest":dest,
        "priority":True
    }
    message_queue_priority[entry].put(vehicule)
    os.kill(os.getppid(), signal.SIGUSR1)
    time.sleep(1)

def lights_priority(vehicule):
    for direction in lights_state:
        lights_state[direction]="RED"
    lights_state[vehicule["entry"]]="GREEN"

def handle_priority_signal(signum, frame):
    global message_queue_priority
    global queue_lock
    print("Priority vehicle detected! Stopping normal traffic.")
    with queue_lock:  # Bloquer l'accès aux ressources partagées
        for direction, msg_queue in message_queue_priority.items():
            while not msg_queue.empty():
                vehicule = msg_queue.get()
                #lights_priority(vehicule)  # Changer les feux pour le véhicule prioritaire
                print(f"Processing priority vehicle {vehicule['id']} from {vehicule['entry']} to {vehicule['dest']}.")


    print("Priority vehicles processed. Resuming normal traffic.")

def display(lights_state,message_queue_normal,message_queue_priority):
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.bind(('localhost', 8080))
    server_socket.listen(5)
    while True:
        try:
            client_socket, client_address = server_socket.accept()
            # Function to get the contents of a multiprocessing.Queue without consuming it
            def peek_queue(queue):
                items = []
                while not queue.empty():
                    item = queue.get()
                    items.append(item)
                # Put items back into the queue
                for item in items:
                    queue.put(item)
                return items

                # Prepare the state information
            state_info = {
                "lights_state": dict(lights_state),  # Current state of the traffic lights
                "vehicles_normal": {direction: peek_queue(queue) for direction, queue in message_queue_normal.items()},  # Normal vehicles in each direction
                "vehicles_priority": {direction: peek_queue(queue) for direction, queue in message_queue_priority.items()}  # Priority vehicles in each direction
            }

                # Send the state information to the client
            try:
                client_socket.sendall(str(state_info).encode('utf-8'))  # Encode and send the state info
            except Exception as e:
                print(f"Error sending data to client: {e}")
            finally:
                client_socket.close()  # Close the client socket after sending data
        except Exception as e:
            print(f"Error in server main loop: {e}")
        time.sleep(1)  # Wait for 1 second before the next update

def turn_left(vehicule,number_destination_reverse):
    # Vérifie si le véhicule doit tourner à gauche
    must_turn_left = (1 + number_destination_reverse[vehicule["entry"]] == number_destination_reverse[vehicule["dest"]]) or \
                     (number_destination_reverse[vehicule["entry"]] == 4 and number_destination_reverse[vehicule["dest"]] == 1)
    return must_turn_left

def duplicate_queue(original_queue):
    new_queue = multiprocessing.Queue()
    while not original_queue.empty():
        item = original_queue.get()
        new_queue.put(item)
        original_queue.put(item)  # Remettre l'élément dans la file d'attente originale
    return new_queue

def put_vehicle_front(queue, vehicle):
    """
    Adds a vehicle to the front of the queue.
    :param queue: The multiprocessing.Queue for the vehicle's entry direction.
    :param vehicle: The vehicle to add to the front of the queue.
    """
    # Create a temporary list to hold the items in the queue
    temp_list = []
    
    # Move all items from the queue to the temporary list
    while not queue.empty():
        temp_list.append(queue.get())
    
    # Add the vehicle to the front of the temporary list
    temp_list.insert(0, vehicle)
    
    # Put all items back into the queue
    for item in temp_list:
        queue.put(item)

def coordinator(message_queue_normal, message_queue_priority, lights_state, number_destination_reverse, opposite, number_destination, id_counter):
    print("Starting coordinator in 20 seconds.")
    time.sleep(5)
    print("Starting coordinator now.")
    signal.signal(signal.SIGUSR1, handle_priority_signal)  # Handle priority signal

    while True:
        print("Coordinator loop running...")
        # Check for priority vehicles first
        for direction, msg_queue in message_queue_priority.items():
            while not msg_queue.empty():
                vehicule = msg_queue.get()
                print(f"Processing priority vehicle {vehicule['id']} from {vehicule['entry']} to {vehicule['dest']}")
                lights_priority(vehicule)  # Change lights for priority vehicle
                print(f"Priority vehicle {vehicule['id']} processed.")

        # Process normal vehicles in a round-robin fashion
        for direction, msg_queue in message_queue_normal.items():
            if not msg_queue.empty():
                vehicule = msg_queue.get()
                print(f"Processing normal vehicle {vehicule['id']} from {vehicule['entry']} to {vehicule['dest']}")
                if lights_state[vehicule["entry"]] == "GREEN":
                    if turn_left(vehicule, number_destination_reverse):
                        with queue_lock:
                            queue = duplicate_queue(message_queue_normal[opposite[vehicule["entry"]]])
                            if queue.empty() or (not queue.empty() and turn_left(queue.get(), number_destination_reverse) == True):
                                print(f"Véhicule {vehicule['id']} tourne à gauche.")
                            else:
                                print(f"Véhicule {vehicule['id']} attend de tourner.")
                                # Put the vehicle back at the front of the queue
                                put_vehicle_front(message_queue_normal[vehicule["entry"]], vehicule)
                    else:
                        print(f"Véhicule {vehicule['id']} part.")
                else:
                    print(f"Véhicule {vehicule['id']} attend.")
                    # Put the vehicle back at the front of the queue
                    put_vehicle_front(message_queue_normal[vehicule["entry"]], vehicule)
                time.sleep(1)

        if all(message_queue_normal[direction].empty() for direction in message_queue_normal):
            print("Toutes les files d'attente sont vides")
            return



if __name__ == "__main__":

    manager = multiprocessing.Manager()
    # Initialize message queues
    message_queue_normal = {
        "N": multiprocessing.Queue(),
        "S": multiprocessing.Queue(),
        "E": multiprocessing.Queue(),
        "W": multiprocessing.Queue()
    }

    # Initialize shared counter
    id_counter = multiprocessing.Value('i', 0)

    # Initialize shared memory
    lights_state =manager.dict({
    "N": "GREEN",
    "S": "GREEN",
    "E": "RED",
    "W": "RED"
    })

    # Initialize opposites
    opposite = manager.dict({
        "N": "S",
        "S": "N",
        "E": "W",
        "W": "E"
    })

    # Initialize number destinations
    number_destination = manager.dict({
        1: "N",
        2: "E",
        3: "S",
        4: "W"
    })

    number_destination_reverse = manager.dict({
        "N": 1,
        "E": 2,
        "S": 3,
        "W": 4
    })

    signal.signal(signal.SIGUSR1, handle_priority_signal)

    kill_switch = multiprocessing.Event()

    # Start processes
    process_light = multiprocessing.Process(target=lights, args=(lights_state,))

    process_normal = multiprocessing.Process(target=normal_traffic_gen, args=(message_queue_normal, number_destination, id_counter))

    process_prio = multiprocessing.Process(target=priority_traffic_gen, args=(message_queue_priority, number_destination, id_counter))

    process_coordinator = multiprocessing.Process(target=coordinator, args=(message_queue_normal, message_queue_priority, lights_state, number_destination_reverse, opposite,number_destination,id_counter))

    process_display = multiprocessing.Process(target=display, args=(lights_state, message_queue_normal, message_queue_priority))

    process_coordinator.start()
    process_light.start()
    process_normal.start()
    process_prio.start()
    process_display.start()

    while True:
   
        time.sleep(3)
        if all(message_queue_normal[direction].empty() for direction in message_queue_normal):
            process_coordinator.terminate()
            process_light.terminate()
            process_normal.terminate()
            process_prio.terminate()
            process_display.terminate()
            break
    
    