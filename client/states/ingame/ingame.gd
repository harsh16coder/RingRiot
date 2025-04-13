extends Node

const packets := preload("res://packets.gd")
const Actor := preload("res://objects/actor/actor.gd")
const Spore := preload("res://objects/spore/spore.gd")
@onready var _line_edit: LineEdit = $UI/LineEdit
@onready var _log: Log = $UI/Log

@onready var _world: Node2D = $World

var _players: Dictionary[int, Actor]
var _spores: Dictionary[int, Spore]


func _ready() -> void:
	WS.connection_closed.connect(_on_ws_connection_closed)
	WS.packet_received.connect(_on_ws_packet_received)

	_line_edit.text_submitted.connect(_on_line_edit_text_entered)

func _on_ws_connection_closed() -> void:
	_log.error("Connection closed")

func _handle_player_msg(sender_id: int, player_msg: packets.PlayerMessage) -> void:
	var actor_id := player_msg.get_id()
	var actor_name := player_msg.get_name()
	var x := player_msg.get_x()
	var y := player_msg.get_y()
	var radius := player_msg.get_radius()
	var speed := player_msg.get_speed()
	
	var is_player := actor_id == GameManager.client_id
	if actor_id not in _players:
		# This is a new player, so we need to create a new actor
		var actor := Actor.instantiate(actor_id, actor_name, x, y, radius, speed, is_player)
		_world.add_child(actor)
		_players[actor_id] = actor
	else:
		# This is an existing player, so we need to update their position
		var actor := _players[actor_id]
		actor.position.x = x
		actor.position.y = y
		var direction := player_msg.get_direction()
		actor.velocity = speed * Vector2.from_angle(direction)

func _on_ws_packet_received(packet: packets.Packet) -> void:
	var sender_id := packet.get_sender_id()
	if packet.has_chat():
		_handle_chat_msg(sender_id, packet.get_chat())
	elif packet.has_player():
		_handle_player_msg(sender_id, packet.get_player())
	elif packet.has_spore():
		_handle_spore_msg(sender_id, packet.get_spore())

func _handle_chat_msg(sender_id: int, chat_msg: packets.ChatMessage) -> void:
	if sender_id in _players:
		var _actor := _players[sender_id]
		_log.chat(_actor.actor_name, chat_msg.get_msg())

func _on_line_edit_text_entered(text: String) -> void:
	var packet := packets.Packet.new()
	var chat_msg := packet.new_chat()
	chat_msg.set_msg(text)
	
	var err = WS.send(packet)
	if err:
		_log.error("Error sending chat message")
	else:
		_log.chat("You", text)
	_line_edit.text = ""

func _handle_spore_msg(sender_id: int, spore_msg: packets.SporeMessage) -> void:
	var spore_id := spore_msg.get_id()
	var x := spore_msg.get_x()
	var y := spore_msg.get_y()
	var radius := spore_msg.get_radius()
	
	if spore_id not in _spores:
		var spore := Spore.instantiate(spore_id, x, y, radius)
		_world.add_child(spore)
		_spores[spore_id] = spore
