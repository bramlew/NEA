<script lang="ts">
    import type { Coords, GeocodeResponse, Response } from "#lib/types.ts"

    const MAX_NO_LOCATIONS = 30 // Maximum no. of locations which can be sent to the API
    const MIN_NO_LOCATIONS = 2 // Minimum no. of locations which can be sent to the API

    // Declare reactive properties
    let type = $state("coords")
    let lon = $state(0)
    let lat = $state(0)
    let loc = $state("")
    let locList: Coords[] = $state([])

    let { initRouteData }: { initRouteData: (data: Response) => void} = $props()

    let enteredLocs: string[] = []

    const updateLocList = async (e: SubmitEvent) => {
        // Update the location list when a new location is submitted
        e.preventDefault()

        switch (type) {
            case "coords":
                // If the location is a set of coordinates, validate whether they are already in the list or not
                const coord: Coords = { lon: lon, lat: lat }
                if (coordsInclude(locList, coord)) {
                    console.error("location already present in list")
                    break
                }
                locList.push(coord)
                break
            case "loc":
                // If the location is a location string, validate whether it is already in the list or not
                const locLower = loc.toLowerCase()
                if (enteredLocs.includes(locLower)) {
                    console.error("location already present in list")
                    break
                }
                enteredLocs.push(locLower)
                // Resolve the coordinates of the location string
                const res = await fetch(`http://localhost:8080/geocode?location=${loc}`)
                const data: GeocodeResponse = await res.json()
                if (!res.ok) {
                    console.error(`error geocoding: ${res.status}: ${res.statusText}`)
                    break
                }
                // Validate that the coordinates of the resolved location are not present in the list or not
                const locCoord: Coords = { lon: data.coord.lon, lat: data.coord.lat }
                if (coordsInclude(locList, locCoord)) {
                    console.error("location already present in list")
                    break
                }
                locList.push(locCoord)
                break
            default:
                console.error("invalid location type")
                break
        }
    }

    const requestOptimise = async () => {
        // Send a request to optimise a route between a list of locations

        if (locList.length > MAX_NO_LOCATIONS || locList.length < MIN_NO_LOCATIONS) {
            // Validate the length of the locations list
            console.error(`no. of locations provided not within valid range, must be between ${MIN_NO_LOCATIONS} and ${MAX_NO_LOCATIONS} inclusive, got ${locList.length}`)
            return
        }

        // Send the request and update the app's data to reflect the response
        const res = await fetch("http://localhost:8080/optimise", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ locations: locList })
        })
        const data: Response = await res.json()
        if (!res.ok) {
            console.error(`error optimising: ${res.status}: ${res.statusText}`)
            return
        }
        initRouteData(data)
    }

    const clearLocs = () => {
        locList = []
        enteredLocs = []
    }

    const coordsInclude = (coords: Coords[], coord: Coords) => {
        // Check if an array of coordinates includes a set of coordinates
        for (let c of coords) {
            if (c.lon === coord.lon && c.lat === coord.lat) {
                return true
            }
        }
        return false
    }
</script>

<div class="flex-1">
    <form class="card mx-auto" onsubmit={updateLocList}>
        <header>
            <h3 class="h3">Location input</h3>
        </header>
        <fieldset class="fieldset">
            <legend class="legend">Enter location type</legend>
            <label class="label">
                <span class="label-text">Coordinates or location</span>
                <label class="label flex">
                    <input class="radio" type="radio" name="locInType" value="coords" bind:group={type} checked />
                    <span class="label-text">Coordinates</span>
                </label>
                <label class="label flex">
                    <input class="radio" type="radio" name="locInType" value="loc" bind:group={type} />
                    <span class="label-text">Location</span>
                </label>
            </label>
        </fieldset>
        <fieldset class="fieldset">
            {#if type === "coords"}
                <legend class="legend">Enter coordinates</legend>
                <label class="label">
                    <span class="label-text">Longitude</span>
                    <input class="input" type="number" placeholder=0 bind:value={lon} min=-180 max=180 step="any" required />
                </label>
                <label class="label">
                    <span class="label-text">Latitude</span>
                    <input class="input" type="number" placeholder=0 bind:value={lat} min=-90 max=90 step="any" required />
                </label>
            {:else if type === "loc"}
                <legend class="legend">Enter location</legend>
                <label class="label">
                    <span class="label-text">Location</span>
                    <input class="input" type="text" placeholder="Eiffel tower" bind:value={loc} required/>
                </label>
            {/if}
        </fieldset>
        <footer class="flex">
            <button type="submit" class="btn preset-filled">Submit</button>
        </footer>
    </form>
    <ul class="ul">
        {#each locList as loc}
            <li class="li">Longitude: {loc.lon} Latitude: {loc.lat}</li>
        {/each}
    </ul>
    <button class="btn preset-filled" onclick={requestOptimise}>Submit request</button>
    <button class="btn preset-filled" onclick={clearLocs}>Clear</button>
</div>
