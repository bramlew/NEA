<script lang="ts">

    import { onDestroy, onMount } from "svelte"
    import Globe, { type GlobeInstance } from "globe.gl"
    import { decode } from "@googlemaps/polyline-codec"
    import type { KeyedLeg } from "#lib/types.ts";

    let { legs } : { legs: KeyedLeg[] } = $props()

    let globeContainer: HTMLElement
    let globe: GlobeInstance
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
    let width = $state(0)
    let height = $state(0)

    onMount(() => {
        globe = new Globe(globeContainer)
            .globeTileEngineUrl((x, y, l) => `https://tiles.versatiles.org/tiles/satellite/${l}/${x}/${y}`)
            .backgroundColor("#091224")
    })

    onDestroy(() => globe._destructor())

    $effect(() => {
        globe.width(width).height(height)
    })

    $effect(() => {
        globe.pathsData(roadRoutesPaths)
            .pathPoints("route")
            .pathLabel("label")
            .pathPointAlt(0)
            .arcsData(airRoutesArcs)
            .arcAltitudeAutoScale(0.3)
            .arcLabel("label")
    })

</script>

<div class="flex-2 w-1/2">
    <div bind:this={globeContainer} bind:clientWidth={width} bind:clientHeight={height} ></div>
</div>
