// Dark neon starfield hero background (raw Three.js, client-only).
// All WebGL work lives here so Svelte components stay thin.
// Callers must dynamic-import this module behind a `browser` guard.
import * as THREE from 'three';

export interface HeroScene {
	destroy(): void;
}

interface Floater {
	mesh: THREE.Mesh;
	spinX: number;
	spinY: number;
	floatSpeed: number;
	floatAmp: number;
	baseY: number;
	phase: number;
}

// Neon-on-dark palette: emerald brand + cyan + violet accents.
const PALETTE = [0x34d399, 0x22d3ee, 0xa78bfa, 0x6ee7b7, 0xe2e8f0];

function makeFloaters(scene: THREE.Scene, rand: () => number): Floater[] {
	const geos: THREE.BufferGeometry[] = [
		new THREE.IcosahedronGeometry(0.55, 0),
		new THREE.OctahedronGeometry(0.6, 0),
		new THREE.TorusGeometry(0.45, 0.14, 10, 26),
		new THREE.BoxGeometry(0.7, 0.7, 0.7),
		new THREE.TetrahedronGeometry(0.6, 0)
	];
	const floaters: Floater[] = [];
	for (let i = 0; i < 12; i++) {
		const color = PALETTE[i % PALETTE.length];
		// Alternate glowing wireframes with dim solids for depth.
		const wireframe = i % 2 === 0;
		const mat = new THREE.MeshStandardMaterial({
			color,
			roughness: 0.4,
			metalness: 0.3,
			wireframe,
			transparent: true,
			opacity: wireframe ? 0.55 : 0.85,
			emissive: color,
			emissiveIntensity: wireframe ? 0.7 : 0.25
		});
		const mesh = new THREE.Mesh(geos[i % geos.length], mat);
		const scale = 0.5 + rand() * 1.1;
		mesh.scale.setScalar(scale);
		mesh.position.set(-7 + rand() * 14, -2.5 + rand() * 5, -4 + rand() * 4.5);
		mesh.rotation.set(rand() * Math.PI, rand() * Math.PI, 0);
		scene.add(mesh);
		floaters.push({
			mesh,
			spinX: (rand() - 0.5) * 0.6,
			spinY: (rand() - 0.5) * 0.8,
			floatSpeed: 0.4 + rand() * 0.8,
			floatAmp: 0.2 + rand() * 0.4,
			baseY: mesh.position.y,
			phase: rand() * Math.PI * 2
		});
	}
	return floaters;
}

function makeStars(scene: THREE.Scene, rand: () => number): THREE.Points {
	const count = 220;
	const pos = new Float32Array(count * 3);
	for (let i = 0; i < count; i++) {
		pos[i * 3] = -10 + rand() * 20;
		pos[i * 3 + 1] = -4 + rand() * 8;
		pos[i * 3 + 2] = -8 + rand() * 6;
	}
	const geo = new THREE.BufferGeometry();
	geo.setAttribute('position', new THREE.BufferAttribute(pos, 3));
	const mat = new THREE.PointsMaterial({
		color: 0x6ee7b7,
		size: 0.045,
		transparent: true,
		opacity: 0.7,
		sizeAttenuation: true
	});
	const stars = new THREE.Points(geo, mat);
	scene.add(stars);
	return stars;
}

export function createHeroScene(canvas: HTMLCanvasElement): HeroScene {
	const renderer = new THREE.WebGLRenderer({ canvas, alpha: true, antialias: true, powerPreference: 'low-power' });
	const scene = new THREE.Scene();
	const camera = new THREE.PerspectiveCamera(55, 1, 0.1, 100);
	camera.position.set(0, 0, 9);

	scene.add(new THREE.AmbientLight(0xffffff, 0.5));
	const key = new THREE.DirectionalLight(0xd1fae5, 1.2);
	key.position.set(4, 6, 6);
	scene.add(key);
	const glowCyan = new THREE.PointLight(0x22d3ee, 14, 30);
	glowCyan.position.set(-6, -2, 3);
	scene.add(glowCyan);
	const glowViolet = new THREE.PointLight(0xa78bfa, 10, 30);
	glowViolet.position.set(6, 3, 2);
	scene.add(glowViolet);

	const floaters = makeFloaters(scene, Math.random);
	const stars = makeStars(scene, Math.random);
	const disposables = new Set<{ dispose(): void }>();
	disposables.add(stars.geometry);
	disposables.add(stars.material as THREE.Material);
	for (const f of floaters) {
		disposables.add(f.mesh.geometry);
		const m = f.mesh.material;
		if (Array.isArray(m)) {
			for (const x of m) disposables.add(x);
		} else {
			disposables.add(m);
		}
	}

	let raf = 0;
	let running = true;
	let visible = true;
	let mouseX = 0;
	let mouseY = 0;
	const clock = new THREE.Clock();

	function resize(): void {
		const w = canvas.clientWidth || 1;
		const h = canvas.clientHeight || 1;
		const cap = window.innerWidth < 720 ? 1.5 : 2;
		renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, cap));
		renderer.setSize(w, h, false);
		camera.aspect = w / h;
		camera.updateProjectionMatrix();
	}

	function frame(): void {
		if (!running) return;
		const t = clock.getElapsedTime();
		for (const f of floaters) {
			f.mesh.rotation.x += f.spinX * 0.016;
			f.mesh.rotation.y += f.spinY * 0.016;
			f.mesh.position.y = f.baseY + Math.sin(t * f.floatSpeed + f.phase) * f.floatAmp;
		}
		stars.rotation.y = t * 0.01;
		camera.position.x += (mouseX * 1.2 - camera.position.x) * 0.04;
		camera.position.y += (mouseY * 0.8 - camera.position.y) * 0.04;
		camera.lookAt(0, 0, 0);
		renderer.render(scene, camera);
		raf = requestAnimationFrame(frame);
	}

	function onMouse(e: MouseEvent): void {
		mouseX = (e.clientX / window.innerWidth - 0.5) * 2;
		mouseY = -(e.clientY / window.innerHeight - 0.5) * 2;
	}

	function onVisibility(): void {
		const active = document.visibilityState === 'visible' && visible;
		if (active && !running) {
			running = true;
			clock.getDelta();
			raf = requestAnimationFrame(frame);
		} else if (!active && running) {
			running = false;
			cancelAnimationFrame(raf);
		}
	}

	const observer = new IntersectionObserver(
		(entries) => {
			visible = entries[0]?.isIntersecting ?? true;
			onVisibility();
		},
		{ threshold: 0 }
	);

	const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
	resize();
	if (reduced) {
		renderer.render(scene, camera);
	} else {
		window.addEventListener('resize', resize);
		window.addEventListener('mousemove', onMouse);
		document.addEventListener('visibilitychange', onVisibility);
		observer.observe(canvas);
		raf = requestAnimationFrame(frame);
	}

	return {
		destroy(): void {
			running = false;
			cancelAnimationFrame(raf);
			window.removeEventListener('resize', resize);
			window.removeEventListener('mousemove', onMouse);
			document.removeEventListener('visibilitychange', onVisibility);
			observer.disconnect();
			for (const d of disposables) d.dispose();
			renderer.dispose();
		}
	};
}
