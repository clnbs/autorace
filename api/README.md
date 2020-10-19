# API definition

## RabbitMQ
### Dynamic server endpoint
##### Player creation 
 - listening route : `autocar.player.creation`
 - accepted data :
```json
{
   "session_uuid":"182a2ed5-d54a-495c-97e2-7e4e5bd806f8",
   "player_name":"toto"
}
```
##### Party creation
 - listening route : `autocar.party.creation`
 - accepted data :
```json
{
   "client_id":"f370a2c8-d6c3-4133-b8a2-6ae6a16e7c5a",
   "seed":321,
   "party_name":"toto party",
   "circuit_config":{
      "seed":321,
      "max_point":250,
      "min_point":50,
      "x_size":4000,
      "y_size":4000
   }
}
```
##### List current parties
 - listening route : `autocar.parties.list`
 - accepted data (Player ID) :
```json
"f370a2c8-d6c3-4133-b8a2-6ae6a16e7c5a"
```

### Dynamic server endpoints
##### Add player to a party
 - listening route : `autocar.party.[@partyID].addPlayer`
 - accepted data :
```json
{
   "session_uuid":"182a2ed5-d54a-495c-97e2-7e4e5bd806f8",
   "player_name":"toto"
}
```

##### Sync
 - listening route : `autocar.party.[@partyID].sync`
 - accepted data : 
```json
{
   "session_uuid":"182a2ed5-d54a-495c-97e2-7e4e5bd806f8",
   "player_name":"toto"
}
```

##### Change game state
 - listening route : `autocar.party.[@partyID].state`
 - accepted data :
```json
{
   "player_token":{
       "session_uuid":"182a2ed5-d54a-495c-97e2-7e4e5bd806f8",
       "player_name":"toto"
   },
   "desired_state": 2
}
```

##### Handle player inputs
 - listening route : `autocar.party.[@partyID].input`
 - accepted data :
```json
{
   "acceleration":0,
   "turning":0,
   "timestamp":"2020-10-15T17:39:04.748300858+02:00",
   "player_uuid":"182a2ed5-d54a-495c-97e2-7e4e5bd806f8"
}
```

### Client endpoints
##### Player creation response
 - listening route : `autocar.player.creation.[@sessionID]`
 - accepted data :
```json
{
   "player_name":"toto",
   "player_uuid":"34a4f55b-03e3-4d0f-ac5a-9e4d240032ca",
   "position":{
      "current_speed":0,
      "current_angle":0,
      "current_position":{
         "x":0,
         "y":0,
         "Angle":0
      }
   },
   "input":{
      "acceleration":0,
      "turning":0,
      "timestamp":"0001-01-01T00:00:00Z",
      "player_uuid":"00000000-0000-0000-0000-000000000000"
   }
}
```

##### Party creation response
 - listening route : `autocar.party.creation.[@clientID]`
 - accepted data
```json
{
   "party_uuid":"0524e4b1-dcf7-4177-b880-af2bcb8363f0",
   "party_name":"testing",
   "map_circuit":{
      "turnpoints":[
         {
            "position":{
               "x":-20.78357379373051,
               "y":332.1042267198367,
               "Angle":0
            }
         },
         {
            "position":{
               "x":-20.58227132252751,
               "y":330.63640667789343,
               "Angle":0
            }
         },
        [......]
         {
            "position":{
               "x":-22.716544160676783,
               "y":324.07827214436367,
               "Angle":0
            }
         },
         {
            "position":{
               "x":-20.78357379373051,
               "y":332.1042267198367,
               "Angle":0
            }
         }
      ]
   },
   "circuit_config":{
      "seed":0,
      "max_point":20,
      "min_point":10,
      "x_size":2000,
      "y_size":2000
   }
}
```

##### List current parties response :
 - listening route : `autocar.parties.list.[@clientID]`
 - accepted data :
```json
[
   "partyID_1",
   "partyID_2",
   "partyID_3",
   "partyID_4"
]
```

##### Receiving a map
 - listening route : `autocar.party.[@partyID].map.[@clientID]`
 - accepted data :
```json
{
   "party_uuid":"0524e4b1-dcf7-4177-b880-af2bcb8363f0",
   "party_name":"testing",
   "map_circuit":{
      "turnpoints":[
         {
            "position":{
               "x":-20.78357379373051,
               "y":332.1042267198367,
               "Angle":0
            }
         },
         {
            "position":{
               "x":-20.58227132252751,
               "y":330.63640667789343,
               "Angle":0
            }
         },
        [......]
         {
            "position":{
               "x":-22.716544160676783,
               "y":324.07827214436367,
               "Angle":0
            }
         },
         {
            "position":{
               "x":-20.78357379373051,
               "y":332.1042267198367,
               "Angle":0
            }
         }
      ]
   },
   "circuit_config":{
      "seed":0,
      "max_point":20,
      "min_point":10,
      "x_size":2000,
      "y_size":2000
   }
}
```

##### Receiving sync message
 - listening route : `autocar.party.[@partyID].sync.[@clientID]`
 - accepted data :
```json
{
   "party_state":1,
   "Competitors":[
      {
         "actor":null,
         "actor_uuid":"7cd39793-0bf7-497c-aeb0-9aa24b68bef0",
         "position":{
            "current_speed":0,
            "current_angle":0,
            "current_position":{
               "x":0,
               "y":0,
               "Angle":0
            }
         }
      },
      {
         "actor":null,
         "actor_uuid":"f56577ed-73bf-40be-8347-662617b6d56e",
         "position":{
            "current_speed":0,
            "current_angle":0,
            "current_position":{
               "x":0,
               "y":0,
               "Angle":0
            }
         }
      },
      {
         "actor":null,
         "actor_uuid":"36f63c86-96ce-4b07-9093-e1f1cf44a13e",
         "position":{
            "current_speed":0,
            "current_angle":0,
            "current_position":{
               "x":0,
               "y":0,
               "Angle":0
            }
         }
      }
   ],
   "MainActor":{
      "actor":null,
      "player":{
         "player_name":"toto",
         "player_uuid":"04a2856e-1fc1-4b61-8f75-f3ecd63ba927",
         "position":{
            "current_speed":0,
            "current_angle":0,
            "current_position":{
               "x":0,
               "y":0,
               "Angle":0
            }
         },
         "input":{
            "acceleration":0,
            "turning":0,
            "timestamp":"0001-01-01T00:00:00Z",
            "player_uuid":"00000000-0000-0000-0000-000000000000"
         }
      }
   }
}
```

##### Receiving a new game state
 - listening route : `autocar.party.[@partyID].sync.[@clientID]`
 - accepted data :
```json
{
   "party_id":"0524e4b1-dcf7-4177-b880-af2bcb8363f0",
   "desired_state":2,
   "new_state":2,
   "message":"OK"
}
```