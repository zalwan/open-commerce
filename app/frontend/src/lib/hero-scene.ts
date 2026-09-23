// Floating-geometry hero background (raw Three.js, client-only).
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

const PALETTE = [0x10b981, 0x0ea5e9, 0xa7f3d0, 0x6ee7b7, 0xffffff];

function makeFloaters(scene: THREE.Scene, rand: () => number): Floater[] {
	const geos: THREE.BufferGeometry[] = [
		new THREE.IcosahedronGeometry(0.55, 0),
		new THREE.OctahedronGeometry(0.6, 0),
		new THREE.TorusGeometry(0.45, 0.16, 12, 28),
		new THREE.BoxGeometry(0.7, 0.7, 0.7),
		new THREE.SphereGeometry(0.4, 20, 16)
	];
	const floaters: Floater[] = [];
	for (let i = 0; i < 14; i++) {
		const color = PALETTE[i % PALETTE.length];
		const wireframe = i % 4 === 3;
		const mat = new THREE.MeshStandardMaterial({
			color,
			roughness: 0.35,
			metalness: 0.15,
			wireframe,
			transparent: true,
			opacity: wireframe ? 0.5 : 0.9
		});
		const mesh = new THREE.Mesh(geos[i % geos.length], mat);
		const scale = 0.6 + rand() * 1.1;
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

export function createHeroScene(canvas: HTMLCanvasElement): HeroScene {
	const renderer = new THREE.WebGLRenderer({ canvas, alpha: true, antialias: true, powerPreference: 'low-power' });
	const scene = new THREE.Scene();
	const camera = new THREE.PerspectiveCamera(55, 1, 0.1, 100);
	camera.position.set(0, 0, 9);

	scene.add(new THREE.AmbientLight(0xffffff, 0.9));
	const key = new THREE.DirectionalLight(0xffffff, 1.6);
	key.position.set(4, 6, 6);
	scene.add(key);
	const rim = new THREE.PointLight(0x0ea5e9, 12, 30);
	rim.position.set(-6, -2, 3);
	scene.add(rim);

	const floaters = makeFloaters(scene, Math.random);
	const geometries = new Set<THREE.BufferGeometry>();
	const materials = new Set<THREE.Material>();
	for (const f of floaters) {
		geometries.add(f.mesh.geometry);
		const m = f.mesh.material;
		if (Array.isArray(m)) {
			for (const x of m) materials.add(x);
		} else {
			materials.add(m);
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
			for (const g of geometries) g.dispose();
			for (const m of materials) m.dispose();
			renderer.dispose();
		}
	};
}
