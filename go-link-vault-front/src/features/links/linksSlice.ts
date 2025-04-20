import { createSlice, createAsyncThunk } from "@reduxjs/toolkit"
import axios from "axios"

export interface Link {
  id: number
  url: string
  title: string
  tags: string[]
}

interface LinksState {
  links: Link[]
  loading: boolean
  error: string | null
}

const initialState: LinksState = {
  links: [],
  loading: false,
  error: null,
}

// Exemple d’appel API
export const fetchLinks = createAsyncThunk("links/fetch", async (_, thunkAPI) => {
  try {
    const response = await axios.get("/links", {
      headers: {
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
    })
    return response.data.data // grâce à notre backend custom response
  } catch (error: any) {
    return thunkAPI.rejectWithValue(error.message)
  }
})

const linksSlice = createSlice({
  name: "links",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchLinks.pending, (state) => {
        state.loading = true
        state.error = null
      })
      .addCase(fetchLinks.fulfilled, (state, action) => {
        state.loading = false
        state.links = action.payload
      })
      .addCase(fetchLinks.rejected, (state, action) => {
        state.loading = false
        state.error = action.payload as string
      })
  },
})

export default linksSlice.reducer

