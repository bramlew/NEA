export type Coords = {
    lon: number
    lat: number
}

export type Response = {
    details: {
        route: Leg[]
        totalWeight: number
    }
}

export type Leg = {
    distance: number
    origin: Coords
    dest: Coords
    isRoad: boolean
    polyline: string
}

export type KeyedLeg = {
    id: number
    details: Leg
}

export type GeocodeResponse = {
    coord: Coords
}
