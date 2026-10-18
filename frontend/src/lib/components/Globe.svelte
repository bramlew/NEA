<script lang="ts">
    import { onDestroy, onMount } from "svelte"
    import Globe, { type GlobeInstance } from "globe.gl"
    import { decode } from "@googlemaps/polyline-codec"
    import type { KeyedLeg } from "#lib/types.ts";

    let { legs } : { legs: KeyedLeg[] } = $props()

    // Declare globe details
    let globeContainer: HTMLElement
    let globe: GlobeInstance

    // Initialise reactive width and height attributes
    let width = $state(0)
    let height = $state(0)

    // Derive route data that can be displayed onto the globe
    let roadRoutes = $derived(legs.filter(leg => leg.details.isRoad))
    let roadRoutesPaths = $derived(roadRoutes.map(leg => { return {
        route: decode(leg.details.polyline),
        label: `Leg ${leg.id}`
    }}))
    let airRoutes = $derived(legs.filter(leg => !leg.details.isRoad))
    let airRoutesArcs = $derived(airRoutes.map(leg => { return {
        startLat: leg.details.origin.lat,
        startLng: leg.details.origin.lon,
        endLat: leg.details.dest.lat,
        endLng: leg.details.dest.lon,
        label: `Leg ${leg.id}`
    }}))

    onMount(() => {
        // Once the page fully loads, set the source of the satellite image data and the background colour to match the
        // rest of the page
        globe = new Globe(globeContainer)
            .globeTileEngineUrl((x, y, l) => `https://tiles.versatiles.org/tiles/satellite/${l}/${x}/${y}`)
            .backgroundColor("#091224")
    })

    onDestroy(() => globe._destructor()) // Destroy the globe if the DOM element is destroyed

    $effect(() => {
        // If the width and height of the container are modified, reflect this change on the globe instance
        globe.width(width).height(height)
    })

    $effect(() => {
        // If new paths/arcs are derived, update the data on the globe to reflect this change
        globe.pathsData(roadRoutesPaths)
            .pathPoints("route")
            .pathLabel("label")
            .pathPointAlt(0.00001)
            .arcsData(airRoutesArcs)
            .arcAltitudeAutoScale(0.3)
            .arcLabel("label")
    })

</script>

<div class="flex-2 w-1/2">
    <div bind:this={globeContainer} bind:clientWidth={width} bind:clientHeight={height} ></div>
</div>
