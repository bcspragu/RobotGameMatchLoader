import numpy as np
from skimage.io import imsave
from skimage.transform import rescale
from sklearn.decomposition import PCA

for x in range(50):
    data = np.loadtxt("rounds/out" + str(x) + ".csv", dtype=int)
    dr = PCA(n_components=5)
    dr.fit(data)
    for i, component in enumerate(dr.components_):
        c = component.reshape(18,18)
        c = rescale(c, 32, order=0)
        imsave("rounds/out" + str(x) + "-component" + str(i) + ".png", c)
