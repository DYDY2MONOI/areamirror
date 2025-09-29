const API_BASE_URL = 'http://localhost:8080'

async function testProfileAPI() {
  console.log('🧪 Test de l\'API Profile...')

  try {
    const response = await fetch(`${API_BASE_URL}/profile`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer test-token'
      }
    })

    console.log('✅ Serveur accessible:', response.status)

    const endpoints = [
      '/profile',
      '/profile/image',
      '/login',
      '/register'
    ]

    for (const endpoint of endpoints) {
      try {
        const testResponse = await fetch(`${API_BASE_URL}${endpoint}`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json'
          }
        })
        console.log(`✅ Endpoint ${endpoint}: ${testResponse.status}`)
      } catch (error) {
        console.log(`❌ Endpoint ${endpoint}: Erreur`)
      }
    }

  } catch (error) {
    console.error('❌ Erreur lors du test:', error)
  }
}

async function testImageUpload() {
  console.log('🧪 Test de l\'upload d\'image...')

  const canvas = document.createElement('canvas')
  canvas.width = 100
  canvas.height = 100
  const ctx = canvas.getContext('2d')

  if (ctx) {
    ctx.fillStyle = '#ff0000'
    ctx.fillRect(0, 0, 100, 100)

    canvas.toBlob(async (blob) => {
      if (blob) {
        const formData = new FormData()
        formData.append('image', blob, 'test-image.png')

        try {
          const response = await fetch(`${API_BASE_URL}/profile/image`, {
            method: 'POST',
            headers: {
              'Authorization': 'Bearer test-token'
            },
            body: formData
          })

          console.log('✅ Upload test:', response.status)
        } catch (error) {
          console.log('❌ Upload test failed:', error)
        }
      }
    }, 'image/png')
  }
}

export { testProfileAPI, testImageUpload }

if (import.meta.env.DEV) {
  console.log('🚀 Démarrage des tests d\'intégration...')
  testProfileAPI()
}
