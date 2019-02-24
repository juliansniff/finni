#version 410 core
out vec4 FragColor;

uniform float opacity;

void main()
{
    FragColor = vec4(0.4, 0.1, 0.3, opacity);
}
