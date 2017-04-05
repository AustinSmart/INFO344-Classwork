# Scripts

### add_droplet_to_load_balancer.sh

>Adds droplet to hardcoded load balancer
>  * Param: 
>    * 1: Droplet id's (CSV)

### create_load_balancer.sh

>Creates a new load balancer
>  * Param: 
>    * 1: Name of new load balancer
  
### create_droplet.sh

>Creates a new VM, starts a docker container, and optionally add to the hardcoded load balancer
>  * Params:   
>    * 1: VM name 
>    * 2: Docker container name - optional
>    * 3: Docker image name - optional
>    * -lb: Add to load balancer - optional
    
### dns.sh

>Returns DNS records of requested type for austinsmart.com
>  *  Params:
>     * 1: Record type (returns all records if not supplied)

### get_droplets.sh

>Returns all droplets for current API token

### get_load_balancer.sh

>Returns the requested load balancer
>   * Param:
>     * 1: Load balancer id
    
### remove_droplets_from_load_balancer.sh

>Gets first load balancer for current API token and prompts for removal/destruction of droplets from it

### verify_ssl.sh

>Returns status of SSL for primary DNS zone of austinsmart.com


