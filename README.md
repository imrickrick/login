# Orbit Labs
Bukid Utility Backend

To Insert or Post a value:

        curl -X POST -H "Content-Type: application/json" 
        -H "X-Farm-Token: ThisIsToken1234567890" 
        -d '{"value":"20"}' http://localhost:8080/api/v1/farm/001
        
To Get the all Recipe:
        
        curl -X GET -H "Content-Type: application/json" 
        -H "X-Farm-Token: ThisIsToken1234567890" 
        http://localhost:8080/api/v1/recipe
        
To Get Some Specific Recipe:
        
        curl -X GET -H "Content-Type: application/json" 
        -H "X-Farm-Token: ThisIsToken1234567890" 
        http://localhost:8080/api/v1/recipe/recipe0
