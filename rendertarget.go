package sf

import (
	"github.com/go-gl-legacy/gl"
)

const vertexCacheSize = 4

type BlendMode uint8

const (
	BlendAlpha    BlendMode = iota // Pixel = Source * Source.a + Dest * (1 - Source.a)
	BlendAdd                       // Pixel = Source + Dest
	BlendMultiply                  // Pixel = Source * Dest
	BlendNone                      /// Pixel = Source
)

// Omg badass render times
type PrimitiveType byte

const (
	Points PrimitiveType = iota
	Lines
	LineStrip
	Triangles
	TriangleStrip
	TriangleFan
	Quads
)

type RenderStates struct {
	BlendMode BlendMode // Blending mode
	Transform Transform // Transform
	Texture   *Texture  // Textures
	//Shader *Shader // Shader
}

type RenderTarget struct {
	size Vector2

	view        *View
	defaultView *View

	// Cache
	glStatesSet    bool      // Are our internal GL states set yet?
	viewChanged    bool      // Has the current view changed since last draw?
	lastBlendMode  BlendMode // Cached blending mode
	lastTextureId  uint64    // Cached texture
	useVertexCache bool      // Did we previously use the vertex cache?
	//vertexCache    [vertexCacheSize]Vertex // Pre-transformed vertices cache

	vpCache [vertexCacheSize]Vector2
	vcCache [vertexCacheSize]Color
	vtCache [vertexCacheSize]Vector2
}

func NewRenderTarget(size Vector2) *RenderTarget {
	rt := &RenderTarget{size: size}
	rt.glStatesSet = false
	rt.defaultView = NewView()
	rt.defaultView.Reset(Rect{0, 0, rt.size.X, rt.size.Y})
	rt.view = NewView()
	*(rt.view) = *(rt.defaultView)
	return rt
}

func (r *RenderTarget) Clear(color Color) {
	gl.ClearColor(gl.GLclampf(float32(color.R)/255), gl.GLclampf(float32(color.G)/255),
		gl.GLclampf(float32(color.B)/255), gl.GLclampf(float32(color.A)/255))
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (r *RenderTarget) SetView(view View) {
	*(r.view) = view
	r.viewChanged = true
}

func (r *RenderTarget) SetSize(size Vector2) {
	r.size = size
	r.defaultView.Reset(Rect{0, 0, r.size.X, r.size.Y})
}

func (r *RenderTarget) Size() Vector2 {
	return r.size
}

func (r *RenderTarget) View() View {
	return *(r.view)
}

func (r *RenderTarget) DefaultView() View {
	return *(r.defaultView)
}

func (r *RenderTarget) Viewport(view *View) Rect {
	w := r.size.X
	h := r.size.Y
	viewport := view.Viewport()

	return Rect{0.5 + w*viewport.Left,
		0.5 + h*viewport.Top,
		w * viewport.W,
		h * viewport.H}
}

func (r *RenderTarget) Render(verts []Vertex, primType PrimitiveType, states RenderStates) {
	// Nothing to draw?
	if len(verts) == 0 {
		return
	}

	// First set the persistent OpenGL states if it's the very first call
	if !r.glStatesSet {
		r.resetGlStates()
	}

	// Check if the vertex count is low enough so that we can pre-transform them
	useVertexCache := len(verts) <= vertexCacheSize
	if useVertexCache {
		// Pre-transform the vertices and store them into the vertex cache
		for i := 0; i < len(verts); i++ {
			r.vpCache[i] = states.Transform.TransformPoint(verts[i].Pos)
			r.vcCache[i] = verts[i].Color
			r.vtCache[i] = verts[i].TexCoords
		}

		// Since vertices are transformed, we must use an identity transform to render them
		if !r.useVertexCache {
			r.applyTransform(IdentityTransform())
		}
	} else {
		r.applyTransform(states.Transform)
	}

	// Apply the view
	if r.viewChanged {
		r.applyCurrentView()
	}

	// Apply the blend mode
	if states.BlendMode != r.lastBlendMode {
		r.applyBlendMode(states.BlendMode)
	}

	// Apply the texture
	var textureId uint64
	if states.Texture != nil {
		textureId = states.Texture.cacheId
	}
	if textureId != r.lastTextureId {
		r.applyTexture(states.Texture)
	}

	// Apply the shader
	// TODO
	/*if states.shader {
		applyShader(states.shader);
	}*/

	// #########################################

	if !useVertexCache {
		// Find the OpenGL primitive type
		modes := [...]gl.GLenum{gl.POINTS, gl.LINES, gl.LINE_STRIP, gl.TRIANGLES,
			gl.TRIANGLE_STRIP, gl.TRIANGLE_FAN, gl.QUADS}
		mode := modes[primType]

		gl.Begin(mode)

		for i, _ := range verts {
			gl.TexCoord2f(verts[i].TexCoords.X, verts[i].TexCoords.Y)
			gl.Color4f(float32(verts[i].Color.R)/255, float32(verts[i].Color.G)/255,
				float32(verts[i].Color.B)/255, float32(verts[i].Color.A)/255)
			gl.Vertex2f(verts[i].Pos.X, verts[i].Pos.Y)
		}

		gl.End()
	}

	// #########################################

	// Setup the pointers to the vertices' components
	// ... and if we already used it previously, we don't need to set the pointers again
	if useVertexCache {
		if !r.useVertexCache {
			gl.VertexPointer(2, gl.FLOAT, 0, r.vpCache[:])
			gl.ColorPointer(4, gl.UNSIGNED_BYTE, 0, r.vcCache[:])
			gl.TexCoordPointer(2, gl.FLOAT, 0, r.vtCache[:])
		}

		// Find the OpenGL primitive type
		modes := [...]gl.GLenum{gl.POINTS, gl.LINES, gl.LINE_STRIP, gl.TRIANGLES,
			gl.TRIANGLE_STRIP, gl.TRIANGLE_FAN, gl.QUADS}
		mode := modes[primType]

		// Draw the primitives
		gl.DrawArrays(mode, 0, len(verts))
	}

	// Unbind the shader, if any
	// TODO
	/*if (states.shader) {
		r.applyShader(nil)
	}*/

	// Update the cache
	r.useVertexCache = useVertexCache
}

func (r *RenderTarget) pushGlStates() {
	gl.PushClientAttrib(gl.CLIENT_ALL_ATTRIB_BITS)
	gl.PushAttrib(gl.ALL_ATTRIB_BITS)
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	gl.MatrixMode(gl.PROJECTION)
	gl.PushMatrix()
	gl.MatrixMode(gl.TEXTURE)
	gl.PushMatrix()

	r.resetGlStates()
}

func (r *RenderTarget) popGlStates() {
	gl.MatrixMode(gl.PROJECTION)
	gl.PopMatrix()
	gl.MatrixMode(gl.MODELVIEW)
	gl.PopMatrix()
	gl.MatrixMode(gl.TEXTURE)
	gl.PopMatrix()
	gl.PopClientAttrib()
	gl.PopAttrib()
}

func (r *RenderTarget) resetGlStates() {
	// Define the default OpenGL states
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.LIGHTING)
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.ALPHA_TEST)
	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.BLEND)
	gl.MatrixMode(gl.MODELVIEW)
	gl.EnableClientState(gl.VERTEX_ARRAY)
	gl.EnableClientState(gl.COLOR_ARRAY)
	gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)
	r.glStatesSet = true

	// Apply the default SFML states
	r.applyBlendMode(BlendAlpha)
	r.applyTransform(IdentityTransform())
	r.applyTexture(nil)
	/*if (Shader::isAvailable()){
		r.applyShader(nil)
	}*/
	r.useVertexCache = false

	// Set the default view
	r.SetView(r.View())
}

func (r *RenderTarget) applyCurrentView() {
	// Set the viewport
	viewport := r.Viewport(r.view)
	top := r.size.Y - (viewport.Top + viewport.H)
	gl.Viewport(int(viewport.Left), int(top), int(viewport.W), int(viewport.H))

	mat := r.view.Transform().Matrix

	// Set the projection matrix
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadMatrixf(&mat)

	// Go back to model-view mode
	gl.MatrixMode(gl.MODELVIEW)

	r.viewChanged = false
}

func (r *RenderTarget) applyBlendMode(mode BlendMode) {
	switch mode {
	// glBlendFuncSeparateEXT is used when available to avoid an incorrect alpha value when the target
	// is a RenderTexture -- in this case the alpha value must be written directly to the target buffer

	// Alpha blending
	default:
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	case BlendAlpha:
		/*if (GLEW_EXT_blend_func_separate) {
		    glBlendFuncSeparateEXT(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA, gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
		} else {*/
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		//}

	// Additive blending
	case BlendAdd:
		/*if GLEW_EXT_blend_func_separate {
			gl.BlendFuncSeparateEXT(gl.SRC_ALPHA, gl.ONE, gl.ONE, gl.ONE)
		} else {*/
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
		//}

	// Multiplicative blending
	case BlendMultiply:
		gl.BlendFunc(gl.DST_COLOR, gl.ZERO)

	// No blending
	case BlendNone:
		gl.BlendFunc(gl.ONE, gl.ZERO)
	}

	r.lastBlendMode = mode
}

func (r *RenderTarget) applyTransform(transform Transform) {
	// No need to call glMatrixMode(gl.MODELVIEW), it is always the
	// current mode (for optimization purpose, since it's the most used)
	gl.LoadMatrixf(&transform.Matrix)
}

func (r *RenderTarget) applyTexture(texture *Texture) {
	texture.Bind(CoordPixels)

	if texture != nil {
		r.lastTextureId = texture.cacheId
	} else {
		r.lastTextureId = 0
	}
}

/*func (r *RenderTarget) applyShader(shader *Shader) {
	Shader::bind(shader);
}*/
