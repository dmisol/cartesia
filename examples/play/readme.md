# Play example

The example uses **cartesia** streaming interface to TTS text to []float32, which is converted to []int16 samples and played using linux **portaudio** library

As a pre-requirement, **portaudio19-dev** is to be installed

The key is expected to be found in conf.json file like 
```
{ 
    "key": "your_cartesia_api_key"
}
```
