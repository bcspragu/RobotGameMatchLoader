import numpy as np
import matplotlib.pyplot as plt

from sklearn.neighbors import KNeighborsClassifier
from sklearn.svm import SVC
from sklearn.tree import DecisionTreeClassifier

from sklearn.cross_validation import cross_val_score

# Data from PCA on bot locations, each round has 10 components
r_data = np.load('train.npy')[:10000]

# Data from robot count by turn
s_data = np.load('synth.npy')[:10000]
# Remove the labels
#val = np.load('validation.npy')

# First round is locations 1 through 10,
# Second is 11 through 20, etc.
def pca_rounds(num_rounds):
    return np.hstack((r_data[:,0].reshape((-1, 1)), r_data[:,range(1,num_rounds*10+1)]))

# First round is locations 1 and 101,
# Second is 2 and 102, etc.
def synth_rounds(num_rounds):
    return np.hstack((s_data[:,0].reshape((-1,1)), s_data[:,range(1,num_rounds+1) + range(1,num_rounds+1)]))


dtc = DecisionTreeClassifier
svc = SVC
knn = KNeighborsClassifier
clfs = [(dtc, "Decision Trees"), (svc, "SVM"), (dtc, "KNN")]
res = [[] for _ in range(len(clfs))]

rounds_to_check = [1,5,10,25,50,75,100]
for i, clf in enumerate(clfs):
    for num_rounds in rounds_to_check:
        rounds = pca_rounds(num_rounds)
        scores = cross_val_score(clf[0](), rounds[:,1:], rounds[:,0], cv=5)
        res[i].append(scores.mean())

fig = plt.figure(1)
fig.suptitle("Cross-Validated Scores by Round (PCA Data)")
for i, result in enumerate(res):
    plt.plot(rounds_to_check, result, label=clfs[i][1])
    plt.xlabel("Number of Rounds Trained On")
    plt.ylabel("Accuracy")

plt.legend()
plt.show()
